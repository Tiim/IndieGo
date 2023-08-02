package activitypub

import (
	"fmt"
	"log"
	"strings"
	"tiim/go-comment-api/config"
)

type activityPubPlugin struct {
	ApPrefix        string `json:"ap_prefix"`
	ActorProfileUrl string `json:"actor_profile_url"`
	ActorName       string `json:"actor_name"`
}

func init() {
	config.RegisterModule(&activityPubPlugin{})
}

func (p *activityPubPlugin) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "activitypub",
		New:  func() config.Module { return new(activityPubPlugin) },
		Docs: config.ConfigDocs{
			DocString: `Minimalist activitypub implementation. TODO: support sending notes when new post is published, store received comments`,
			Fields: map[string]string{
				"ApPrefix":        "Prefix of the url to the activitypub endpoint. Example: 'https://comments.tiim.ch'. The string '/ap' is automatically appendend.",
				"ActorProfileUrl": "Url to the profile of the actor, example 'https://tiim.ch'",
				"ActorName":       "Name of the actor. Example 'user'",
			},
		},
	}
}

func (p *activityPubPlugin) Load(c config.GlobalConfig, _ interface{}, logger *log.Logger) (config.ModuleInstance, error) {

	if p.ApPrefix == "" || p.ActorName == "" {
		return nil, fmt.Errorf("ap_prefix or actor_name is not specified")
	}

	p.ApPrefix = strings.TrimSuffix(p.ApPrefix, "/")
	var apModule config.ApiPluginInstance = &activityPubModule{
        store: apStore {
            baseUrl: p.ApPrefix,
            actorProfileUrl: p.ActorProfileUrl,
            actorName: p.ActorName,
        },
		apPrefix:        p.ApPrefix,
		actorProfileUrl: p.ActorProfileUrl,
		actorName:       p.ActorName,
	}
	return apModule, nil
}
