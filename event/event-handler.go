package event

import (
	"log"
	"tiim/go-comment-api/model"
)

type HandlerList struct {
	handlers []Handler
}

func NewHandlerList(handlers []Handler) *HandlerList {
	return &HandlerList{handlers: handlers}
}

func (l *HandlerList) OnNewComment(c *model.GenericComment) (bool, error) {
	for _, h := range l.handlers {
		if ok, err := h.OnNewComment(c); !ok || err != nil {
			log.Printf("error in event handler %s (OnNewComment): %s", h.Name(), err)
			return ok, err
		}
	}
	return true, nil
}

func (l *HandlerList) OnDeleteComment(c *model.GenericComment) (bool, error) {
	for _, h := range l.handlers {
		if ok, err := h.OnDeleteComment(c); !ok || err != nil {
			log.Printf("error in event handler %s (OnDeleteComment): %s", h.Name(), err)
			return ok, err
		}
	}
	return true, nil
}

func (l *HandlerList) Name() string {
	return "HandlerList"
}
