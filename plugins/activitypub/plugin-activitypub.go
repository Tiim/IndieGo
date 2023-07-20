package activitypub

import (
	"log"
	"tiim/go-comment-api/config"
)

type activityPubPlugin struct{}

func init() {
	config.RegisterModule(&activityPubPlugin{})
}

func (p *activityPubPlugin) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "activitypub",
		New:  func() config.Module { return new(activityPubPlugin) },
		Docs: config.ConfigDocs{
			DocString: `Minimalist activitypub implementation. TODO: support sending notes when new post is published, store received comments`,
		},
	}
}

func (p *activityPubPlugin) Load(config config.GlobalConfig, _ interface{}, logger *log.Logger) (config.ModuleInstance, error) {

	return newActivityPubModule(), nil
}
