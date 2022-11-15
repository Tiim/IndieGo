package indieauth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"time"
)

type AuthCode struct {
	code                string
	clientId            string
	redirectUri         string
	scope               string
	state               string
	codeChallenge       string
	codeChallengeMethod string
	me                  string
	ts                  time.Time
}

type Store interface {
	StoreAuthCode(authCode *AuthCode) error
	GetAuthCode(code string) (*AuthCode, error)
	DeleteAuthCode(code string) error
}

type SQLiteStore struct {
	db                *sql.DB
	authCodeValidTime time.Duration
}

func NewSQLiteStore(db *sql.DB, authCodeValidTime time.Duration) *SQLiteStore {
	return &SQLiteStore{db: db, authCodeValidTime: authCodeValidTime}
}

func newAuthCode(redirectUri, clientId, scope, state, codeChallenge, codeChallengeMethod, me string) (*AuthCode, error) {
	buffer := make([]byte, 32)
	_, err := rand.Read(buffer)
	if err != nil {
		return nil, err
	}
	code := base64.RawURLEncoding.EncodeToString(buffer)
	return &AuthCode{
		code:                code,
		redirectUri:         redirectUri,
		clientId:            clientId,
		scope:               scope,
		state:               state,
		codeChallenge:       codeChallenge,
		codeChallengeMethod: codeChallengeMethod,
		me:                  me,
		ts:                  time.Now(),
	}, nil
}

func (s *SQLiteStore) StoreAuthCode(authCode *AuthCode) error {
	_, err := s.db.Exec("INSERT INTO indieauth_auth_codes (code, client_id, redirect_uri, scope, state, code_challenge, code_challenge_method, me, ts) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		authCode.code, authCode.clientId, authCode.redirectUri, authCode.scope, authCode.state, authCode.codeChallenge, authCode.codeChallengeMethod, authCode.me, authCode.ts.Format(time.RFC3339))
	return err
}

func (s *SQLiteStore) GetAuthCode(code string) (*AuthCode, error) {
	var authCode AuthCode
	var ts string
	row := s.db.QueryRow("SELECT code, client_id, redirect_uri, scope, state, code_challenge, code_challenge_method, me, ts FROM indieauth_auth_codes WHERE code = ? AND ts > ?",
		code, time.Now().Add(-s.authCodeValidTime).Format(time.RFC3339))
	err := row.Scan(
		&authCode.code,
		&authCode.clientId,
		&authCode.redirectUri,
		&authCode.scope,
		&authCode.state,
		&authCode.codeChallenge,
		&authCode.codeChallengeMethod,
		&authCode.me,
		&ts,
	)
	if err != nil {
		return nil, err
	}

	authCode.ts, err = time.Parse(time.RFC3339, ts)
	if err != nil {
		return nil, err
	}

	return &authCode, nil
}

func (s *SQLiteStore) DeleteAuthCode(code string) error {
	_, err := s.db.Exec("DELETE FROM indieauth_auth_codes WHERE code = ?", code)
	return err
}

func (s *SQLiteStore) CleanUp() error {
	_, err := s.db.Exec("DELETE FROM indieauth_auth_codes WHERE ts < ?", time.Now().Add(-s.authCodeValidTime).Format(time.RFC3339))
	return err
}
