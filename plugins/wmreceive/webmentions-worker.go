package wmrecv

import (
	"log"
	"time"
)

type mentionsQueueWorker struct {
	store   webmentionsStore
	checker *WebmentionChecker
	logger  *log.Logger
}

func newMentionsQueueWorker(store webmentionsStore, checker *WebmentionChecker, logger *log.Logger) *mentionsQueueWorker {
	worker := &mentionsQueueWorker{store: store, checker: checker, logger: logger}
	go worker.run()
	return worker
}

func (w *mentionsQueueWorker) run() {
	for {
		time.Sleep(1 * time.Second)
		wm, err := w.store.NextWebmentionFromQueue()
		if err != nil {
			w.logger.Println(err)
			return
		}
		err = w.processNextWebmention(wm)
		if err != nil {
			w.logger.Printf("Error processing webmention: %s", err)
		}
	}
}

func (w *mentionsQueueWorker) processNextWebmention(wm *QueuedWebmention) error {
	err := w.checker.CheckWebmentionValid(wm.webmention)
	if err != nil {
		w.logger.Printf("Webmention %s failed checks: %v", wm.webmention.Source, err)
		return w.store.MarkInvalid(wm, err.Error())
	} else {
		w.logger.Printf("Webmention %s passed checks", wm.webmention.Source)
		return w.store.MarkSuccess(wm)
	}

}
