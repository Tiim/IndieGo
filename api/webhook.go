package api

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebhookCallback func(*gin.Context) error

type webhookItem struct {
	name      string
	callback  WebhookCallback
	validator WebhookValidator
}

type webhookApiModule struct {
	webhooks []webhookItem
}

func NewWebhookModule() *webhookApiModule {
	return &webhookApiModule{webhooks: make([]webhookItem, 0)}
}

func (w *webhookApiModule) RegisterWebhook(name string, callback WebhookCallback, validator WebhookValidator) {
	w.webhooks = append(w.webhooks, webhookItem{name: name, callback: callback, validator: validator})
}

func (w *webhookApiModule) Name() string {
	return "Webhook"
}

func (w *webhookApiModule) Init(r *gin.Engine, templates fs.FS) error {
	return nil
}

func (w *webhookApiModule) RegisterRoutes(r *gin.Engine) error {
	for _, webhook := range w.webhooks {
		r.POST("/webhook/"+webhook.name, func(c *gin.Context) {
			if err := webhook.validator(c); err != nil {
				log.Printf("Webhook %s: invalid request: %v", webhook.name, err)
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			err := webhook.callback(c)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		})
	}
	return nil
}
