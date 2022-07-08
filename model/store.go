package model

import "time"

type Store interface {
	NewComment(c *Comment) error
	GetAllComments(since time.Time) ([]Comment, error)
	GetCommentsForPost(page string, since time.Time) ([]Comment, error)
	DeleteComment(id string) error
	GetComment(id string) (Comment, error)
}