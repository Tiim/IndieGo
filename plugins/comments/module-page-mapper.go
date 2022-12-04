package comments

import (
	"encoding/json"
	"tiim/go-comment-api/config"
)

type pageMapperModule struct{}
type pageMapperModuleData struct {
	Format string `json:"format"`
}

func init() {
	config.RegisterModule(&pageMapperModule{})
}

func (m *pageMapperModule) Name() string {
	return "comments-page-mapper"
}

func (m *pageMapperModule) Load(data json.RawMessage, config config.GlobalConfig) (config.ModuleInstance, error) {
	d := pageMapperModuleData{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}
	return &formatPageMapper{format: d.Format}, nil
}
