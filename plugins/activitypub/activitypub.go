package activitypub

import (
	"fmt"
	"tiim/go-comment-api/config"

	"git.sr.ht/~mariusor/lw"
	"github.com/gin-gonic/gin"
	"github.com/go-ap/activitypub"
	"github.com/go-ap/webfinger"
)

type activityPubModule struct {
	group *gin.RouterGroup
}

func newActivityPubModule() *activityPubModule {
	return &activityPubModule{}
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
	store := &webfingerActorStore{}
	wf := webfinger.New(logger, store)
    if c.Query("resource") == "" {
        c.AbortWithError(404, fmt.Errorf("No parameter 'resource' given"))
        return
    }
	wf.HandleWebFinger(c.Writer, c.Request)
}

type webfingerActorStore struct{}

func (d *webfingerActorStore) Load(iri activitypub.IRI) (activitypub.Item, error) {
	fmt.Printf("iri: %s\n", iri)
	actor := activitypub.ActorNew(
		"https://tiim.ch/",
		activitypub.ActorType,
	)

    actor.PreferredUsername.Set(activitypub.DefaultLang, activitypub.Content("user"))
    actor.ID = activitypub.IRI("https://comments.tiim.ch/ap/users/1111")
    actor.URL = actor.ID
    

	return actor, nil
}

// TODO:
// required routes:
// - webfinger
// - actor
// - inbox -> receives "mentions" and post answers
// - outbox -> lists blog posts
// - following -> return only actor specified in config
// - followed -> return actors from the database
