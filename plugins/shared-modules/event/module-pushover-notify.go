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
		Docs: config.ConfigDocs{
			DocString: `Pushover notification module. This module sends a Pushover notification when a new comment is submitted.`,
			Fields: map[string]string{
				"ApiToken": "Pushover API token. You can find it on the settings page of the app on your pushover account.",
				"UserKey":  "Pushover user key. The personal key for the pushover account you want to send the notification to.",
			},
		},
	}
}

func (m *pushoverNotifyModule) Load(config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {

	return newPushoverNotify(m.ApiToken, m.UserKey, *config.HttpClient), nil
}
