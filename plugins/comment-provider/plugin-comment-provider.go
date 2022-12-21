package commentprovider

import (
	"fmt"
	"log"
	"tiim/go-comment-api/config"
)

type CommentProviderPlugin struct{}

func init() {
	config.RegisterModule(&CommentProviderPlugin{})
	config.RegisterInterface("comment-provider.provider")
}

func (p *CommentProviderPlugin) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "comment-provider",
		New:  func() config.Module { return new(CommentProviderPlugin) },
		Docs: config.ConfigDocs{
			DocString: `Comment provider plugin. This plugin serves comments from the /comment endpoint. 
				The comments are provided by previously loaded plugins that register a comment-provider.provider interface.`,
			Fields: map[string]string{},
		},
	}
}

func (p *CommentProviderPlugin) Load(c config.GlobalConfig, _ interface{}) (config.ModuleInstance, error) {
	providerInter := c.Config.GetInterfaces("comment-provider.provider")
	providers := make([]CommentProvider, len(providerInter))
	for i, iface := range providerInter {
		p, ok := iface.(CommentProvider)
		if !ok {
			return nil, fmt.Errorf("interface is not a CommentProvider: %T", iface)
		}
		providers[i] = p
	}

	log.Printf("Loaded %d comment providers", len(providers))

	var providerModule config.ApiPluginInstance = newCommentProviderModule(providers)
	return providerModule, nil
}
