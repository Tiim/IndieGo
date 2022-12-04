package admin

import (
	"encoding/json"
	"tiim/go-comment-api/config"
)

type AdminPlugin struct{}
type AdminPluginData struct {
	Password string `json:"password"`
}

func init() {
	config.RegisterPlugin(&AdminPlugin{})
}

func (p *AdminPlugin) Name() string {
	return "admin"
}

func (p *AdminPlugin) Load(data json.RawMessage, config config.GlobalConfig) (config.PluginInstance, error) {
	d := AdminPluginData{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}
	return newAdminModule(d.Password), nil
}
