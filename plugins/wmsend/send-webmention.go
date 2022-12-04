package wmsend

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/mmcdole/gofeed"
	"willnorris.com/go/webmention"
)

type wmSend struct {
	store     WmSendStore
	rss       string
	client    *http.Client
	scheduler *gocron.Scheduler
	interval  time.Duration
}

type FeedItem struct {
	uid     string
	link    string
	content string
	updated *time.Time
	baseUrl string
}

func newWmSend(store WmSendStore, client *http.Client, rss string, scheduler *gocron.Scheduler, interval time.Duration) *wmSend {
	return &wmSend{
		store:     store,
		rss:       rss,
		client:    client,
		scheduler: scheduler,
		interval:  interval,
	}
}

func (w *wmSend) Name() string {
	return "wmsend"
}

func (w *wmSend) Init() error {
	return nil
}

func (w *wmSend) Start() error {
	w.scheduler.Every(w.interval).Do(w.SendNow)
	return nil
}

func (w *wmSend) SendNow() {
	go func() {
		err := w.doFetchAndSend()
		if err != nil {
			log.Printf("unable to send webmentions: %v", err)
		}
	}()
}

func (w *wmSend) doFetchAndSend() error {

	log.Println("Sending webmentions...")

	feed, err := w.getFeedItems()
	if err != nil {
		return fmt.Errorf("unable to get feed items: %v", err)
	}
	for _, item := range feed {
		updated, err := w.store.IsItemUpdated(item)
		if err != nil {
			log.Printf("unable to check if item is updated: %v", err)
			continue
		}
		if updated {
			err := w.sendWebmentions(item)
			if err != nil {
				log.Printf("unable to send webmentions: %v", err)
				continue
			}
		}
	}
	return nil
}

func (w *wmSend) getFeedItems() ([]FeedItem, error) {
	fp := gofeed.NewParser()

	resp, err := w.client.Get(w.rss)

	if err != nil {
		return nil, err
	}

	feed, err := fp.Parse(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("unable to parse feed: %v", err)
	}

	items := make([]FeedItem, len(feed.Items))

	for i, item := range feed.Items {
		time := item.UpdatedParsed
		if time == nil {
			time = item.PublishedParsed
		}
		items[i] = FeedItem{uid: item.Link, content: item.Content, updated: time, link: item.Link, baseUrl: feed.Link}
	}

	return items, nil
}

func (w *wmSend) sendWebmentions(item FeedItem) error {
	reader := strings.NewReader(item.content)
	links, err := webmention.DiscoverLinksFromReader(reader, item.baseUrl, "")

	if err != nil {
		return fmt.Errorf("unable to discover links: %w", err)
	}

	savedLinks, err := w.store.GetUrlsForFeedItem(item)

	if err != nil {
		return fmt.Errorf("unable to get saved links: %w", err)
	}

	links = append(links, savedLinks...)

	linksSet := make(map[string]struct{})
	for _, link := range links {
		linksSet[link] = struct{}{}
	}

	err = w.store.SetUrlsForFeedItem(item, links)

	if err != nil {
		return fmt.Errorf("unable to save links: %w", err)
	}

	wmClient := webmention.New(w.client)

	for link := range linksSet {
		endpoint, err := wmClient.DiscoverEndpoint(link)
		if err != nil {
			log.Printf("unable to discover endpoint for url %s: %v", link, err)
		} else {
			log.Printf("sending webmention from %s to %s", item.link, link)
			wmClient.SendWebmention(endpoint, item.link, link)
		}
	}

	return nil
}
