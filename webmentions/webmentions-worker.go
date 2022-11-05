package webmentions

import (
	"log"
	"time"
)

type mentionsQueueWorker struct {
	store *webmentionsStore
}

func NewMentionsQueueWorker(store *webmentionsStore) *mentionsQueueWorker {
	worker := &mentionsQueueWorker{store: store}
	go worker.run()
	return worker
}

func (w *mentionsQueueWorker) run() {
	for {
		time.Sleep(1 * time.Second)
		wm, err := w.store.NextWebmentionFromQueue()
		if err != nil {
			log.Println(err)
			return
		}
		w.processNextWebmention(wm)
	}
}

func (w *mentionsQueueWorker) processNextWebmention(wm *QueuedWebmention) error {
	err := checkWebmentionValid(wm.webmention)
	if err != nil {
		return w.store.MarkInvalid(wm, err.Error())
	} else {
		return w.store.MarkSuccess(wm)
	}

}
