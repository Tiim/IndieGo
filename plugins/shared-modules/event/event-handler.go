package event

import (
	"log"
	"tiim/go-comment-api/model"
)

type handlerList struct {
	handlers []Handler
	logger   *log.Logger
}

func (l *handlerList) OnNewComment(c *model.GenericComment) (bool, error) {
	for _, h := range l.handlers {
		if ok, err := h.OnNewComment(c); !ok || err != nil {
			l.logger.Printf("error in event handler %s (OnNewComment): %s", h.Name(), err)
			return ok, err
		}
	}
	return true, nil
}

func (l *handlerList) OnDeleteComment(c *model.GenericComment) (bool, error) {
	for _, h := range l.handlers {
		if ok, err := h.OnDeleteComment(c); !ok || err != nil {
			l.logger.Printf("error in event handler %s (OnDeleteComment): %s", h.Name(), err)
			return ok, err
		}
	}
	return true, nil
}

func (l *handlerList) Name() string {
	return "HandlerList"
}
