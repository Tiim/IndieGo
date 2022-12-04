package indieauth

import (
	"crypto/sha256"
	"encoding/base64"
)

var challenges = map[string]codeChallenge{
	"S256":  &challengeS256{},
	"plain": &challengePlain{},
}

type codeChallenge interface {
	Name() string
	Verify(codeVerifier, codeChallenge string) bool
}

type challengeS256 struct {
}

func (s *challengeS256) Name() string {
	return "S256"
}

func (s *challengeS256) Verify(codeVerifier, codeChallenge string) bool {
	sha256 := sha256.New()
	sha256.Write([]byte(codeVerifier))
	return codeChallenge == base64.RawURLEncoding.EncodeToString(sha256.Sum(nil))
}

type challengePlain struct {
}

func (s *challengePlain) Name() string {
	return "plain"
}

func (s *challengePlain) Verify(codeVerifier, codeChallenge string) bool {
	return codeVerifier == "" && codeChallenge == ""
}
