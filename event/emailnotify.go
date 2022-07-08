package event

import (
	"fmt"
	"log"
	"net/smtp"
	"tiim/go-comment-api/model"
)

type EmailNotify struct {
	From     string
	To       string
	Subject  string
	Password string
	SmtpHost string
	SmtpPort string
}

func (n *EmailNotify) OnNewComment(c *model.Comment) (bool, error) {
	go n.doSendEmail(c)
	return true, nil
}

func (n *EmailNotify) doSendEmail(c *model.Comment) {
	auth := smtp.PlainAuth("", n.From, n.Password, n.SmtpHost)
	text := fmt.Sprintf("New Comment\nid:\t%s\nfrom:\t%s <%s>\npage:\t%s\n\n%s", c.Id, c.Name, c.Email, c.Page, c.Content)
	data := fmt.Sprintf("Content-Type: text/plain; charset=UTF-8\nSubject: %s\nFrom: Comment System <%s>\nReply-To: %s <%s>\n\n%s\n", n.Subject, n.From, c.Name, c.Email, text)
	err := smtp.SendMail(n.SmtpHost+":"+n.SmtpPort, auth, n.From, []string{n.To}, []byte(data))
	if err != nil {
		log.Printf("Error sending notification email: %s", err)
	}
	log.Printf("notification email sent to %s", n.To)
}

func (n *EmailNotify) OnDeleteComment(c *model.Comment) (bool, error) {
	return true, nil
}

func (n *EmailNotify) Name() string {
	return "EmailNotify"
}
