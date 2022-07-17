package model

import (
	"io"
	"time"
)

type Store interface {
	NewComment(c *Comment) error
	GetAllComments(since time.Time) ([]Comment, error)
	GetCommentsForPost(page string, since time.Time) ([]Comment, error)
	DeleteComment(id string) error
	GetComment(id string) (Comment, error)
}

type SubscribtionStore interface {
	Store
	Unsubscribe(secret string) (Comment, error)
	UnsubscribeAll(email string) ([]Comment, error)
}

type CleanupStore interface {
	CleanUp() error
}

type BackupStore interface {
	Backup() (io.Reader, error)
}
