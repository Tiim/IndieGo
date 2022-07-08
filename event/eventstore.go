package event

import (
	"fmt"
	"os"
	"tiim/go-comment-api/model"
)

type Store struct {
	store         model.Store
	eventhandlers []Handler
}

func NewEventStore(store model.Store, eventhandlers []Handler) *Store {
	return &Store{store: store, eventhandlers: eventhandlers}
}

func (s *Store) NewComment(c *model.Comment) error {
	for _, h := range s.eventhandlers {
		next, err := h.OnNewComment(c)

		if err != nil {
			fmt.Fprintf(os.Stderr, "[%s] on new comment %s: %s", h.Name(), c.Id, err)
		}
		if !next {
			return nil
		}
	}
	return s.store.NewComment(c)
}

func (s *Store) GetAllComments() ([]model.Comment, error) {
	return s.store.GetAllComments()
}

func (s *Store) GetCommentsForPost(page string) ([]model.Comment, error) {
	return s.store.GetCommentsForPost(page)
}

func (s *Store) DeleteComment(id string) error {
	comment, err := s.store.GetComment(id)

	if err != nil {
		return fmt.Errorf("error getting comment %s: %w", id, err)
	}

	for _, h := range s.eventhandlers {
		next, err := h.OnDeleteComment(&comment)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[%s] on delete comment %s: %s", h.Name(), id, err)
		}
		if !next {
			return nil
		}
	}
	return s.store.DeleteComment(id)
}

func (s *Store) GetComment(id string) (model.Comment, error) {
	return s.store.GetComment(id)
}
