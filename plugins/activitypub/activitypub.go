package activitypub

import (
	"fmt"
	"net/url"
	"tiim/go-comment-api/config"

	"git.sr.ht/~mariusor/lw"
	"github.com/gin-gonic/gin"
	"github.com/go-ap/activitypub"
	"github.com/go-ap/webfinger"
)

type activityPubModule struct {
	apPrefix        string
	actorProfileUrl string
	actorName       string
	group           *gin.RouterGroup
}

func (m *activityPubModule) Name() string {
	return "activitypub"
}

func (m *activityPubModule) Init(config config.GlobalConfig) error {
	return nil
}

func (m *activityPubModule) Start() error {
	return nil
}

func (m *activityPubModule) InitGroups(r *gin.Engine) error {
	m.group = r.Group("/ap")
	return nil
}

func (m *activityPubModule) RegisterRoutes(r *gin.Engine) error {
	r.GET(".well-known/webfinger", m.handleWebfinger)
	return nil
}

func (m *activityPubModule) handleWebfinger(c *gin.Context) {
	logger := lw.Dev()
	// TODO: don't instantiate this on every call
	baseUrl, err := url.Parse(m.apPrefix)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	store := &webfingerActorStore{
		baseUrl:         m.apPrefix,
		actorProfileUrl: m.actorProfileUrl,
		actorName:       m.actorName,
		host:            baseUrl.Hostname(),
	}
	wf := webfinger.New(logger, store)
	if c.Query("resource") == "" {
		c.AbortWithError(404, fmt.Errorf("No parameter 'resource' given"))
		return
	}
	wf.HandleWebFinger(c.Writer, c.Request)
}

// TODO:
// required routes:
// - webfinger
// - actor
// - inbox -> receives "mentions" and post answers
// - outbox -> lists blog posts
// - following -> return only actor specified in config
// - followed -> return actors from the database
