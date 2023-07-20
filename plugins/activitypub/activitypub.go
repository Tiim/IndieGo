package activitypub

import "github.com/gin-gonic/gin"


type activityPubModule struct {
    

}

func newActivityPubModule() *activityPubModule {
    return &activityPubModule {}
}

func (m *activityPubModule) Name() string {
    return "activitypub"
}

func (m *activityPubModule) Start() error {
    return nil
}

func (m *activityPubModule) InitGroups(r *gin.Engine) error {
    return nil
}

func (m *activityPubModule) RegisterRoutes(r *gin.Engine) error {
    
    return nil
}

// TODO:
// required routes: 
// - webfinger
// - actor 
// - inbox -> receives "mentions" and post answers
// - outbox -> lists blog posts
// - following -> return only actor specified in config
// - followed -> return actors from the database




