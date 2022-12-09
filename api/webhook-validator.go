package api

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type WebhookValidator func(*gin.Context) error

var DefaultWebhookValidator WebhookValidator = func(c *gin.Context) error { return nil }

func NewGithubValidator(key string) WebhookValidator {
	return func(c *gin.Context) error {
		if !isValidSignature(c.Request, key) {
			return errors.New("invalid signature")
		}
		return nil
	}
}

func isValidSignature(r *http.Request, key string) bool {
	// Assuming a non-empty header
	gotHash := strings.SplitN(r.Header.Get("X-Hub-Signature"), "=", 2)
	if gotHash[0] != "sha1" {
		return false
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Cannot read the request body: %s\n", err)
		return false
	}

	hash := hmac.New(sha1.New, []byte(key))
	if _, err := hash.Write(b); err != nil {
		log.Printf("Cannot compute the HMAC for request: %s\n", err)
		return false
	}

	expectedHash := hex.EncodeToString(hash.Sum(nil))
	return gotHash[1] == expectedHash
}
