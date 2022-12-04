package event

import (
	"fmt"
	"log"
	"net/smtp"
	"tiim/go-comment-api/model"

	"github.com/jordan-wright/email"
)

type emailNotify struct {
	from     string
	to       string
	subject  string
	username string
	password string
	smtpHost string
	smtpPort string
}

func (n *emailNotify) OnNewComment(c *model.GenericComment) (bool, error) {
	go n.doSendEmail(c)
	return true, nil
}

func (n *emailNotify) doSendEmail(c *model.GenericComment) {

	log.Printf("sending notification email from %s to %s", n.from, n.to)

	e := email.NewEmail()
	e.From = n.from
	e.To = []string{n.to}
	e.Subject = n.subject
	e.Text = []byte(fmt.Sprintf("New %s\nid:\t%s\nfrom:\t%s <%s>\npage:\t%s\n\n%s", c.Type, c.Id, c.Name, c.FromEmail, c.Page, c.Content))

	log.Printf("sending mail: %s:%s user:%s", n.smtpHost, n.smtpPort, n.username)
	err := e.Send(n.smtpHost+":"+n.smtpPort, smtp.PlainAuth("", n.username, n.password, n.smtpHost))

	if err != nil {
		log.Println("error sending notification email:", err)
	} else {
		log.Println("notification email sent")
	}
}

func (n *emailNotify) OnDeleteComment(c *model.GenericComment) (bool, error) {
	return true, nil
}

func (n *emailNotify) Name() string {
	return "EmailNotify"
}
