package indieauth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"
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

type AccessToken struct {
	scope     string
	clientId  string
	issuedAt  time.Time
	expiresAt time.Time
}

type Store interface {
	StoreAuthCode(authCode *AuthCode) error
	GetAuthCode(code string) (*AuthCode, error)
	DeleteAuthCode(code string) error
	UpdateScope(code, scope string) error
	RedeemAccessToken(authCode string) (*AccessToken, error)
}

type sQLiteStore struct {
	db                 *sql.DB
	authCodeValidTime  time.Duration
	authTokenValidTime time.Duration
	logger             *log.Logger
}

func NewSQLiteStore(db *sql.DB, authCodeValidTime, authTokenValidTime time.Duration, logger *log.Logger) *sQLiteStore {
	return &sQLiteStore{db: db, authCodeValidTime: authCodeValidTime, authTokenValidTime: authTokenValidTime, logger: logger}
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

func (s *sQLiteStore) StoreAuthCode(authCode *AuthCode) error {
	_, err := s.db.Exec("INSERT INTO indieauth_auth_codes (code, client_id, redirect_uri, scope, state, code_challenge, code_challenge_method, me, ts) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		authCode.code, authCode.clientId, authCode.redirectUri, authCode.scope, authCode.state, authCode.codeChallenge, authCode.codeChallengeMethod, authCode.me, authCode.ts.Format(time.RFC3339))
	return err
}

func (s *sQLiteStore) GetAuthCode(code string) (*AuthCode, error) {
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

func (s *sQLiteStore) DeleteAuthCode(code string) error {
	_, err := s.db.Exec("DELETE FROM indieauth_auth_codes WHERE code = ?", code)
	return err
}

func (s *sQLiteStore) UpdateScope(code, scope string) error {
	_, err := s.db.Exec("UPDATE indieauth_auth_codes SET scope = ? WHERE code = ?", scope, code)
	return err
}

func (s *sQLiteStore) CleanUp() error {
	_, err := s.db.Exec("DELETE FROM indieauth_auth_codes WHERE ts < ?", time.Now().Add(-s.authCodeValidTime).Format(time.RFC3339))
	return err
}

func (s *sQLiteStore) RedeemAccessToken(authCode string) (*AccessToken, error) {
	tx, err := s.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	var ts string
	var authCodeObj AuthCode
	row := tx.QueryRow("SELECT code, client_id, scope, me, ts FROM indieauth_auth_codes WHERE code = ? AND ts > ?",
		authCode, time.Now().Add(-s.authCodeValidTime).Format(time.RFC3339))
	err = row.Scan(
		&authCodeObj.code,
		&authCodeObj.clientId,
		&authCodeObj.scope,
		&authCodeObj.me,
		&ts,
	)
	if err != nil {
		return nil, err
	}
	authCodeObj.ts, _ = time.Parse(time.RFC3339, ts)

	_, err = tx.Exec("DELETE FROM indieauth_auth_codes WHERE code = ?", authCode)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &AccessToken{
		scope:     authCodeObj.scope,
		clientId:  authCodeObj.clientId,
		issuedAt:  time.Now(),
		expiresAt: time.Now().Add(s.authTokenValidTime),
	}, nil
}
