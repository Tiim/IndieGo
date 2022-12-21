package event

import (
	"tiim/go-comment-api/config"
)

type emailNotifyModule struct {
	From     string `json:"email_from"`
	To       string `json:"email_to"`
	Subject  string `json:"email_subject"`
	Username string `json:"username"`
	Password string `json:"password"`
	SmtpHost string `json:"smtp_host"`
	SmtpPort string `json:"smtp_port"`
}

func init() {
	config.RegisterModule(&emailNotifyModule{})
}

func (m *emailNotifyModule) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "event.mention.email-notify",
		New:  func() config.Module { return new(emailNotifyModule) },
		Docs: config.ConfigDocs{
			DocString: `Email notification module. This module sends an email when a new comment is submitted.`,
			Fields: map[string]string{
				"From":     "Email address to send from.",
				"To":       "Email address to send to.",
				"Subject":  "Email subject.",
				"Username": "Username for SMTP server.",
				"Password": "Password for SMTP server.",
				"SmtpHost": "SMTP host.",
				"SmtpPort": "SMTP port.",
			},
		},
	}
}

func (m *emailNotifyModule) Load(config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	return &emailNotify{
		from:     m.From,
		to:       m.To,
		subject:  m.Subject,
		username: m.Username,
		password: m.Password,
		smtpHost: m.SmtpHost,
		smtpPort: m.SmtpPort,
	}, nil
}
