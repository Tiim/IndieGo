package event

import (
	"tiim/go-comment-api/config"
)

type pushoverNotifyModule struct {
	ApiToken string `json:"api_token"`
	UserKey  string `json:"user_key"`
}

func init() {
	config.RegisterModule(&pushoverNotifyModule{})
}

func (m *pushoverNotifyModule) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "event.mention.pushover-notify",
		New:  func() config.Module { return new(pushoverNotifyModule) },
	}
}

func (m *pushoverNotifyModule) Load(config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {

	return newPushoverNotify(m.ApiToken, m.UserKey, *config.HttpClient), nil
}
