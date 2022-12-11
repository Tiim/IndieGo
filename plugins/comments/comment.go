package comments

import "tiim/go-comment-api/model"

type comment struct {
	Id                string `json:"id"`
	ReplyTo           string `json:"reply_to"`
	Timestamp         string `json:"timestamp"`
	Page              string `json:"page"`
	Url               string `json:"url"`
	Content           string `json:"content"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	Notify            bool   `json:"notify"`
	UnsubscribeSecret string `json:"-"`
}

func (c *comment) ToGenericComment() model.GenericComment {
	return model.GenericComment{
		Id:        c.Id,
		Type:      "comment",
		ReplyTo:   c.ReplyTo,
		FromEmail: c.Email,
		Timestamp: c.Timestamp,
		Page:      c.Page,
		Url:       c.Url,
		Content:   c.Content,
		Name:      c.Name,
	}
}
