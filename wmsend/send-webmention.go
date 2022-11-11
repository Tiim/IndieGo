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
	store  WmSendStore
	rss    string
	client *http.Client
}

type FeedItem struct {
	uid     string
	content string
	updated *time.Time
	baseUrl string
}

func NewWmSend(store WmSendStore, client *http.Client, rss string) *wmSend {
	return &wmSend{store: store, client: client, rss: rss}
}

func (w *wmSend) Start() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Hour().Do(w.doFetchAndSend)
	w.doFetchAndSend()
}

func (w *wmSend) doFetchAndSend() error {

	log.Println("Sending webmentions...")

	feed, err := w.getFeedItems()
	if err != nil {
		log.Printf("unable to get feed items: %v", err)
		return err
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
		items[i] = FeedItem{uid: item.Link, content: item.Content, updated: time, baseUrl: feed.Link}
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

	w.store.SetUrlsForFeedItem(item, links)

	wmClient := webmention.New(w.client)

	for link := range linksSet {
		endpoint, err := wmClient.DiscoverEndpoint(link)
		if err != nil {
			log.Printf("unable to discover endpoint for url %s: %v", link, err)
		} else {
			log.Printf("sending webmention from to %s", link)
			wmClient.SendWebmention(endpoint, item.baseUrl, link)
		}
	}

	return nil
}
