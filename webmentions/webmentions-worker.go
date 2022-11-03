package webmentions

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type mentionsQueueWorker struct {
	store *webmentionsStore
	next  chan bool
}

func NewMentionsQueueWorker(store *webmentionsStore) *mentionsQueueWorker {
	channel := make(chan bool)
	worker := &mentionsQueueWorker{store: store, next: channel}
	go worker.run()
	return worker
}

func (w *mentionsQueueWorker) Ping() {
	select {
	case w.next <- true:
	default:
	}
}

func (w *mentionsQueueWorker) run() {
	fmt.Println("Webmentions worker started")
	for {
		for {
			time.Sleep(1 * time.Second)
			wm, err := w.store.getNextWebmentionFromQueue()
			if err != nil || wm == nil {
				if err != nil {
					fmt.Println("Error getting next webmention from queue:", err)
				}
				break
			}
			w.processNextWebmention(wm)
		}
		<-w.next
	}
}

func (w *mentionsQueueWorker) processNextWebmention(wm *QueuedWebmention) {
	fmt.Printf("Processing webmention %s\n", wm)
	passed, err := checkWebmentionHTML(wm)
	if err != nil {
		fmt.Println("Error processing webmention:", err)
		w.store.updateNextTry(wm)
		return
	}

	fmt.Printf("passed: %t\n", passed)

	if passed {
		w.store.moveWebmentionFromQueueToProcessed(wm)
	} else {
		w.store.deleteFromQueue(wm)
	}

}

func checkWebmentionHTML(wm *QueuedWebmention) (bool, error) {
	log.Println("Checking webmention", wm)
	res, err := http.Get(wm.webmention.Source)
	if err != nil {
		return false, fmt.Errorf("error fetching source url: %w", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return false, fmt.Errorf("error parsing source html: %w", err)
	}

	targetUrl, err := url.ParseRequestURI(wm.webmention.Target)
	if err != nil {
		return false, fmt.Errorf("error parsing target url: %w", err)
	}

	foundTarget := false

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if foundTarget {
			return
		}
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		foundUrl, err := url.ParseRequestURI(href)
		if err == nil && *foundUrl == *targetUrl {
			foundTarget = true
		}
	})

	return foundTarget, nil
}
