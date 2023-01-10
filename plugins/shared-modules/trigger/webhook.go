package trigger

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"tiim/go-comment-api/config"

	"github.com/gin-gonic/gin"
)

type webhookModule struct {
	name      string
	callback  []Callback
	validator WebhookValidator
	logger    *log.Logger
}

func newWebhookModule(name string, logger *log.Logger) *webhookModule {
	return &webhookModule{name: name, validator: DefaultWebhookValidator, logger: logger}
}

func (w *webhookModule) AddCallback(callback Callback) {
	w.callback = append(w.callback, callback)
}

func (w *webhookModule) Name() string {
	return "Webhook"
}

func (w *webhookModule) Start() error {
	return nil
}

func (w *webhookModule) Init(config.GlobalConfig) error {
	return nil
}

func (w *webhookModule) RegisterRoutes(r *gin.Engine) error {
	r.POST("/webhook/"+w.name, func(c *gin.Context) {
		if err := w.validator(c); err != nil {
			w.logger.Printf("Webhook %s: invalid request: %v", w.name, err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		for _, callback := range w.callback {
			callback()
		}
	})
	return nil
}

func (w *webhookModule) SetValidator(v WebhookValidator) {
	w.validator = v
}

type WebhookValidator func(*gin.Context) error

var DefaultWebhookValidator WebhookValidator = func(c *gin.Context) error { return nil }

func newGithubValidator(key string, logger *log.Logger) WebhookValidator {
	return func(c *gin.Context) error {
		if !isValidSignature(c.Request, key, logger) {
			return errors.New("invalid signature")
		}
		return nil
	}
}

func isValidSignature(r *http.Request, key string, logger *log.Logger) bool {
	// Assuming a non-empty header
	gotHash := strings.SplitN(r.Header.Get("X-Hub-Signature"), "=", 2)
	if gotHash[0] != "sha1" {
		return false
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Printf("Cannot read the request body: %s\n", err)
		return false
	}

	hash := hmac.New(sha1.New, []byte(key))
	if _, err := hash.Write(b); err != nil {
		logger.Printf("Cannot compute the HMAC for request: %s\n", err)
		return false
	}

	expectedHash := hex.EncodeToString(hash.Sum(nil))
	return gotHash[1] == expectedHash
}
