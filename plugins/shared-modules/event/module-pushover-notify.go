package event

import (
	"encoding/json"
	"tiim/go-comment-api/config"
)

type PushoverNotifyModule struct{}
type PushoverNotifyModuleData struct {
	ApiToken string `json:"api_token"`
	UserKey  string `json:"user_key"`
}

func init() {
	config.RegisterModule(&PushoverNotifyModule{})
}

func (m *PushoverNotifyModule) Name() string {
	return "event-pushover-notify"
}

func (m *PushoverNotifyModule) Load(data json.RawMessage, config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	var d PushoverNotifyModuleData
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}
	return newPushoverNotify(d.ApiToken, d.UserKey, *config.HttpClient), nil
}
