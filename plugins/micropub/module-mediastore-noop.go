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
	}
}

func (p *MediastoreNoopModule) Load(config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	return nopMediaStore{}, nil
}
