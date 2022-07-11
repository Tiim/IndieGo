package event

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
	"tiim/go-comment-api/model"

	"github.com/jordan-wright/email"
)

type replyEmail struct {
	store    model.Store
	from     string
	subject  string
	username string
	password string
	smtpHost string
	smtpPort string
	baseUrl  string
	template *template.Template
}

func NewReplyEmail(store model.Store, from, subject, username, password, smtpHost, smtpPort, baseUrl string) *replyEmail {

	html := `
	<html>
		<body>
			<p>
				<b> {{ .NewComment.Name }} </b> replied to your comment:
			</p>
			<blockquote>
				<p>From: <b>{{ .YourComment.Name }}</b> (You)</p>
				<p>{{ .YourComment.Content }}</p>
			</blockquote>
			<blockquote>
				<p>From: <b>{{ .NewComment.Name }}</b></p>
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

func (n *replyEmail) OnNewComment(c *model.Comment) (bool, error) {
	if c.ReplyTo != "" {
		go n.doSendEmail(*c)
	}
	return true, nil
}

func (n *replyEmail) doSendEmail(c model.Comment) {
	commentChain := n.collectReplyChain(c)

	for _, comment := range commentChain {
		log.Printf("sending reply notification email from %s to %s", n.from, comment.Email)

		var buf bytes.Buffer
		err := n.template.Execute(&buf, struct {
			NewComment  model.Comment
			YourComment model.Comment
			BaseUrl     string
		}{
			NewComment:  c,
			YourComment: comment,
			BaseUrl:     n.baseUrl,
		})

		if err != nil {
			log.Println("error sending reply notification email:", err)
			continue
		}

		e := email.NewEmail()
		e.From = n.from
		e.To = []string{string(comment.Email)}
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

func (n *replyEmail) collectReplyChain(currentComment model.Comment) []model.Comment {
	comments := make([]model.Comment, 0)
	for currentComment.ReplyTo != "" {
		var err error
		currentComment, err = n.store.GetComment(currentComment.ReplyTo)
		if err != nil {
			log.Printf("error getting comment #%d: %s", len(comments), err)
		} else {
			comments = append(comments, currentComment)
		}
	}
	return comments
}

func (n *replyEmail) OnDeleteComment(c *model.Comment) (bool, error) {
	return true, nil
}

func (n *replyEmail) Name() string {
	return "ReplyEmail"
}
