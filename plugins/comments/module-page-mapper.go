package comments

import (
	"log"
	"tiim/go-comment-api/config"
)

type pageMapperModule struct {
	Format string `json:"format"`
}

func init() {
	config.RegisterModule(&pageMapperModule{})
}

func (p *pageMapperModule) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "comments.page-mapper.format",
		New:  func() config.Module { return new(pageMapperModule) },
		Docs: config.ConfigDocs{
			DocString: `Page mapper module. This module is responsible for formatting comment urls.`,
			Fields: map[string]string{
				"Format": "The format string to use for mapping comment ids to urls. The format can contain the following placeholders: {id} - The comment id, {page} - The page the comment is on. Example: https://example.com/{page}#comment-{id}",
			},
		},
	}
}

func (p *pageMapperModule) Load(config config.GlobalConfig, args interface{}, logger *log.Logger) (config.ModuleInstance, error) {
	return &formatPageMapper{format: p.Format, logger: logger}, nil
}
