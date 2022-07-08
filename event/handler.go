package event

import "tiim/go-comment-api/model"

type Handler interface {
	OnNewComment(c *model.Comment) (bool, error)
	OnDeleteComment(c *model.Comment) (bool, error)
	Name() string
}
