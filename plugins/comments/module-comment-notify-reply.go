package comments

import (
	"encoding/json"
	"fmt"
	"tiim/go-comment-api/config"
)

// can only be used as a child of a comment store
type commentNotifyReplyModule struct{}
type commentNotifyReplyModuleData struct {
	EmailFrom    string `json:"email_from"`
	EmailSubject string `json:"email_subject"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	SmtpHost     string `json:"smtp_host"`
	SmtpPort     string `json:"smtp_port"`
	BaseUrl      string `json:"base_url"`
}

func init() {
	config.RegisterModule(&commentNotifyReplyModule{})
}

func (p *commentNotifyReplyModule) Name() string {
	return "comment-notify-reply"
}

func (p *commentNotifyReplyModule) Load(data json.RawMessage, config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	commentStore, ok := args.(commentStore)
	if !ok {
		return nil, fmt.Errorf("can only be used as a child of a comments.commentStore module")
	}

	d := commentNotifyReplyModuleData{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}
	return newReplyEmail(
		commentStore,
		d.EmailFrom,
		d.EmailSubject,
		d.Username,
		d.Password,
		d.SmtpHost,
		d.SmtpPort,
		d.BaseUrl,
	), nil
}
