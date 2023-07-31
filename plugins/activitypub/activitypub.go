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

type webfingerActorStore struct {
	baseUrl         string
	actorProfileUrl string
	actorName       string
	host            string
}

func (d *webfingerActorStore) Load(iri activitypub.IRI) (activitypub.Item, error) {
	fmt.Printf("iri: %s\n", iri)
	url, err := iri.URL()
	if err != nil {
		return nil, err
	}
	if url.Path == "/" && url.Hostname() == d.host {
        fmt.Println("Service actor")
		return d.buildServiceActor()
	} else {
        fmt.Printf("path %s, host: %s ref %s\n", url.Path, url.Hostname(), d.host)
        fmt.Println("User actor")
		return d.buildPersonActor()
	}

}

func (d *webfingerActorStore) buildServiceActor() (activitypub.Item, error) {
	actor := activitypub.ActorNew(
		activitypub.IRI(d.baseUrl+"/ap"),
		activitypub.ServiceType,
	)

	actor.PreferredUsername.Set(activitypub.DefaultLang, activitypub.Content("IndieGo Server"))
	actor.URL = actor.ID

	actor.Name.Set(activitypub.DefaultLang, activitypub.Content("indiego"))

	return actor, nil
}

func (d *webfingerActorStore) buildPersonActor() (activitypub.Item, error) {
	actor := activitypub.ActorNew(
		activitypub.IRI(d.baseUrl+"/ap/users/"+d.actorName),
		activitypub.PersonType,
	)

	actor.PreferredUsername.Set(activitypub.DefaultLang, activitypub.Content(d.actorName))
	if d.actorProfileUrl != "" {
		actor.URL = activitypub.IRI(d.actorProfileUrl)
	}
	actor.Name.Set(activitypub.DefaultLang, activitypub.Content(d.actorName))

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
