package event

import (
	"encoding/json"
	"tiim/go-comment-api/config"
)

type emailNotifyModule struct{}
type emailNotifyModuleData struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Username string `json:"username"`
	Password string `json:"password"`
	SmtpHost string `json:"smtp_host"`
	SmtpPort string `json:"smtp_port"`
}

func init() {
	config.RegisterModule(&emailNotifyModule{})
}

func (m *emailNotifyModule) Name() string {
	return "event-email-notify"
}

func (m *emailNotifyModule) Load(data json.RawMessage, config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	var d emailNotifyModuleData
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}
	return &emailNotify{
		from:     d.From,
		to:       d.To,
		subject:  d.Subject,
		username: d.Username,
		password: d.Password,
		smtpHost: d.SmtpHost,
		smtpPort: d.SmtpPort,
	}, nil
}
