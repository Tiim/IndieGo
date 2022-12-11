package commentprovider

import (
	"tiim/go-comment-api/model"
	"time"
)

type CommentProvider interface {
	GetGenericCommentsForPage(page string, since time.Time) ([]model.GenericComment, error)
	GetAllGenericComments(since time.Time) ([]model.GenericComment, error)
}
