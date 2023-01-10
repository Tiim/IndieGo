package comments

import (
	"fmt"
	"log"
	"tiim/go-comment-api/config"
)

// can only be used as a child of a comment store
type commentNotifyReplyModule struct {
	// The email address to use as the sender.
	EmailFrom string `json:"email_from"`
	// The subject to use for the email.
	EmailSubject string `json:"email_subject"`
	// The username for the smtp server.
	Username string `json:"username"`
	// The password for the smtp server.
	Password string `json:"password"`
	// The hostname of the smtp server.
	SmtpHost string `json:"smtp_host"`
	// The port of the smtp server.
	SmtpPort string `json:"smtp_port"`
	// The base url of the website. Used to generate the link to the comment.
	BaseUrl string `json:"base_url"`
}

func init() {
	config.RegisterModule(&commentNotifyReplyModule{})
}

func (p *commentNotifyReplyModule) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "event.mention.email-reply",
		New:  func() config.Module { return new(commentNotifyReplyModule) },
		Docs: config.ConfigDocs{
			DocString: `Email reply notification module. This module sends an email to a commenter when a new reply to their comment chain is submitted.`,
			Fields: map[string]string{
				"EmailFrom":    "The email address to use as the sender.",
				"EmailSubject": "The subject to use for the email.",
				"Username":     "The username for the smtp server.",
				"Password":     "The password for the smtp server.",
				"SmtpHost":     "The hostname of the smtp server.",
				"SmtpPort":     "The port of the smtp server.",
				"BaseUrl":      "The base url of the website. Used to generate the link to the comment.",
			},
		},
	}
}

func (p *commentNotifyReplyModule) Load(config config.GlobalConfig, args interface{}, logger *log.Logger) (config.ModuleInstance, error) {
	commentStore, ok := args.(commentStore)
	if !ok {
		return nil, fmt.Errorf("can only be used as a child of a comments.store module")
	}

	return newReplyEmail(
		commentStore,
		p.EmailFrom,
		p.EmailSubject,
		p.Username,
		p.Password,
		p.SmtpHost,
		p.SmtpPort,
		p.BaseUrl,
		logger,
	), nil
}
