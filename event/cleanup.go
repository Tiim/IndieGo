package event

import (
	"tiim/go-comment-api/model"
	"time"
)

type CleanUp struct {
	Store       model.CleanupStore
	lastCleanup *time.Time
}

func (s *CleanUp) OnNewComment(c *model.GenericComment) (bool, error) {
	t := time.Now()
	if s.lastCleanup == nil || t.Sub(*s.lastCleanup) > time.Hour*24 {
		s.lastCleanup = &t
		go s.Store.CleanUp()
	}
	return true, nil
}
func (s *CleanUp) OnDeleteComment(c *model.GenericComment) (bool, error) {
	t := time.Now()
	if s.lastCleanup == nil || t.Sub(*s.lastCleanup) > time.Hour*24 {
		s.lastCleanup = &t
		go s.Store.CleanUp()
	}
	return true, nil
}
func (s *CleanUp) Name() string {
	return "CleanUp"
}
