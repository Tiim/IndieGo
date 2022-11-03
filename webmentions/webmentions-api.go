package webmentions

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type webmentionsModule struct {
	store  *webmentionsStore
	worker *mentionsQueueWorker
}

func NewApi(store *webmentionsStore, worker *mentionsQueueWorker) *webmentionsModule {
	im := webmentionsModule{store: store, worker: worker}
	return &im
}

func (ui *webmentionsModule) Name() string {
	return "Webmentions"
}

func (ui *webmentionsModule) Init(r *gin.Engine) error {
	return nil
}

func (ui *webmentionsModule) RegisterRoutes(r *gin.Engine) error {
	r.POST("/wm/webmentions", ui.handlePostWebmention)
	return nil
}

func (ui *webmentionsModule) handlePostWebmention(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if !c.Request.Form.Has("source") || !c.Request.Form.Has("target") {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing source or target"))
		return
	}

	source := c.Request.Form.Get("source")
	target := c.Request.Form.Get("target")

	wm, err := newWebmention(source, target)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := ui.store.scheduleForProcessing(wm); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ui.worker.Ping()

	c.Status(http.StatusAccepted)
}
