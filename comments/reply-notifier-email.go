package comments

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"tiim/go-comment-api/model"

	"github.com/jordan-wright/email"
)

type replyEmail struct {
	store    *commentStore
	from     string
	subject  string
	username string
	password string
	smtpHost string
	smtpPort string
	baseUrl  string
	template *template.Template
}

func NewReplyEmail(store *commentStore, from, subject, username, password,
	smtpHost, smtpPort, baseUrl string) *replyEmail {

	html := `
	<html>
		<body>
			<p>
				<b> {{ .NewComment.Name }} </b> replied to your comment:
			</p>
			<blockquote>
				<p>From: <a href="{{ .YourCommentUrl }}"><b>{{ .YourComment.Name }}</b> (You)</a></p>
				<p>{{ .YourComment.Content }}</p>
			</blockquote>
			<blockquote>
				<p>From: <a href="{{ .NewCommentUrl }}"><b>{{ .NewComment.Name }}</b></a></p>
				<p>{{ .NewComment.Content }}</p>
			</blockquote>
			<p>
				<a href="{{ .BaseUrl }}/unsubscribe/comment/{{ .YourComment.UnsubscribeSecret }}">Unsubscribe</a>
			</p>
		</body>
	</html>	
	`

	template := template.Must(template.New("replyEmail").Parse(html))

	return &replyEmail{
		store:    store,
		from:     from,
		subject:  subject,
		username: username,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		baseUrl:  baseUrl,
		template: template,
	}
}

func (n *replyEmail) OnNewComment(c *model.GenericComment) (bool, error) {
	if c.ReplyTo != "" {
		go n.doSendEmail(*c)
	}
	return true, nil
}

func (n *replyEmail) doSendEmail(c model.GenericComment) {
	commentChain, err := n.collectReplyChain(c)

	if err != nil {
		log.Println("error collecting reply chain for email notifications", err)
	}

	for _, cChain := range commentChain {
		log.Printf("sending reply notification email from %s to %s", n.from, cChain.Email)

		var buf bytes.Buffer
		err := n.template.Execute(&buf, struct {
			NewComment     model.GenericComment
			YourComment    comment
			BaseUrl        string
			YourCommentUrl string
			NewCommentUrl  string
		}{
			NewComment:     c,
			YourComment:    cChain,
			BaseUrl:        n.baseUrl,
			YourCommentUrl: cChain.Url,
			NewCommentUrl:  c.Page,
		})

		if err != nil {
			log.Println("error sending reply notification email:", err)
			continue
		}

		e := email.NewEmail()
		e.From = n.from
		e.To = []string{string(cChain.Email)}
		e.Subject = n.subject
		e.HTML = buf.Bytes()

		log.Printf("sending mail: %s:%s user:%s", n.smtpHost, n.smtpPort, n.username)
		err = e.Send(n.smtpHost+":"+n.smtpPort, smtp.PlainAuth("", n.username, n.password, n.smtpHost))

		if err != nil {
			log.Println("error sending notification email:", err)
		} else {
			log.Println("notification email sent")
		}
	}
}

func (n *replyEmail) collectReplyChain(currentComment model.GenericComment) ([]comment, error) {
	comments := make([]comment, 0)
	replyTo := currentComment.ReplyTo
	for replyTo != "" {
		var err error
		cmt, err := n.store.GetComment(currentComment.ReplyTo, nil)
		if err != nil {
			return nil, fmt.Errorf("error getting comment #%d: %s", len(comments), err)
		}
		replyTo = cmt.ReplyTo
	}
	return comments, nil
}

func (n *replyEmail) OnDeleteComment(c *model.GenericComment) (bool, error) {
	return true, nil
}

func (n *replyEmail) Name() string {
	return "ReplyEmail"
}
