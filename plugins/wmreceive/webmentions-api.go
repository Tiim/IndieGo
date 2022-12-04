package wmrecv

import (
	"fmt"
	"net/http"
	"strings"
	"tiim/go-comment-api/config"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

type webmentionsModule struct {
	store     webmentionsStore
	worker    *mentionsQueueWorker
	scheduler *gocron.Scheduler
}

func newApi(store webmentionsStore, worker *mentionsQueueWorker, scheduler *gocron.Scheduler) *webmentionsModule {
	im := webmentionsModule{store: store, worker: worker, scheduler: scheduler}
	return &im
}

func (ui *webmentionsModule) Name() string {
	return "webmention-receive"
}

func (ui *webmentionsModule) Init(config config.GlobalConfig) error {
	return nil
}

func (ui *webmentionsModule) InitGroups(r *gin.Engine) error {
	return nil
}

func (ui *webmentionsModule) RegisterRoutes(r *gin.Engine) error {
	r.POST("/wm/webmentions", ui.handlePostWebmention)
	return nil
}

func (ui *webmentionsModule) Start() error {
	ui.scheduler.Every(4).Hours().Do(ui.store.RefetchQueue)
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

	wm, err := NewWebmention(source, target)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := ui.store.ScheduleForProcessing(wm); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// check if the referer header is set and if it is, redirect back to it
	if referer := c.Request.Header.Get("Referer"); referer != "" && strings.HasSuffix(referer, "/admin") {
		c.Redirect(http.StatusFound, referer)
		return
	}

	c.Status(http.StatusAccepted)
}
