package event

import "tiim/go-comment-api/model"

type Handler interface {
	Name() string
	OnNewComment(c *model.GenericComment) (bool, error)
	OnDeleteComment(c *model.GenericComment) (bool, error)
}
