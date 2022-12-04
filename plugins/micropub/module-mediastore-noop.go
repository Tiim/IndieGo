package micropub

import (
	"encoding/json"
	"tiim/go-comment-api/config"
)

type MediastoreNoopModule struct{}

func init() {
	config.RegisterModule(&MediastoreNoopModule{})
}

func (p *MediastoreNoopModule) Name() string {
	return "micropub-mediastore-noop"
}

func (p *MediastoreNoopModule) Load(data json.RawMessage, config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	return nopMediaStore{}, nil
}
