package comments

import (
	"tiim/go-comment-api/config"
)

type pageMapperModule struct {
	// The format string to use for mapping comment ids to urls.
	// The format can contain the following placeholders:
	// "{id}" - The comment id
	// "{page}" - The page the comment is on
	// Example: "https://example.com/{page}#comment-{id}"
	Format string `json:"format"`
}

func init() {
	config.RegisterModule(&pageMapperModule{})
}

func (p *pageMapperModule) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "comments.page-mapper.format",
		New:  func() config.Module { return new(pageMapperModule) },
	}
}

func (p *pageMapperModule) Load(config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	return &formatPageMapper{format: p.Format}, nil
}
