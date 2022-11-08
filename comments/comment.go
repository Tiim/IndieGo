package comments

import "tiim/go-comment-api/model"

type comment struct {
	Id                string
	ReplyTo           string
	Timestamp         string
	Page              string
	Url               string
	Content           string
	Name              string
	Email             string
	Notify            bool
	UnsubscribeSecret string
}

func (c *comment) ToGenericComment() model.GenericComment {
	return model.GenericComment{
		Id:        c.Id,
		Type:      "comment",
		ReplyTo:   c.ReplyTo,
		FromEmail: c.Email,
		Timestamp: c.Timestamp,
		Page:      c.Url,
		Content:   c.Content,
		Name:      c.Name,
	}
}
