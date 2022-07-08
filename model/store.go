package model

type Store interface {
	NewComment(c *Comment) error
	GetAllComments() ([]Comment, error)
	GetCommentsForPost(page string) ([]Comment, error)
	DeleteComment(id string) error
	GetComment(id string) (Comment, error)
}
