package micropub

import (
	"tiim/go-comment-api/config"
)

type MediastoreNoopModule struct{}

func init() {
	config.RegisterModule(&MediastoreNoopModule{})
}

func (p *MediastoreNoopModule) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "micropub.media-store.noop",
		New:  func() config.Module { return new(MediastoreNoopModule) },
		Docs: config.ConfigDocs{
			DocString: `Noop media store module. This media store discards all media. Useful for testing or when you don't want any media to be stored.`,
		},
	}
}

func (p *MediastoreNoopModule) Load(config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	return nopMediaStore{}, nil
}
