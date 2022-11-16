package wmsend

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

type TestWMStore struct {
	urls []string
}

func (t *TestWMStore) IsItemUpdated(item FeedItem) (bool, error) {
	return true, nil
}

func (t *TestWMStore) GetUrlsForFeedItem(item FeedItem) ([]string, error) {
	return t.urls, nil
}

func (t *TestWMStore) SetUrlsForFeedItem(item FeedItem, urls []string) error {
	return nil
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func Test_wmSend_getFeedItems(t *testing.T) {

	client := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.String() == "https://tiim.ch/blog/rss.xml" {
			f, err := os.Open("testdata/rss/tiim.ch.rss.xml")
			if err != nil {
				t.Errorf("unable to open testdata: %v", err)
			}
			return &http.Response{
				Body: f,
				Header: http.Header{
					"Content-Type": []string{"application/rss+xml"},
				},
			}
		}
		t.Errorf("unexpected request: %v", req.URL.String())
		return nil
	})

	wmsend := NewWmSend(&TestWMStore{}, client, "https://tiim.ch/blog/rss.xml")
	feed, err := wmsend.getFeedItems()

	if err != nil {
		t.Fatal(err)
	}

	if len(feed) != 1 {
		t.Errorf("expected 1 feed item, got %d", len(feed))
		return
	}
	if feed[0].uid != "https://tiim.ch/blog/2022-09-27-sveltekit-ssr-with-urql" {
		t.Errorf("expected uid to be https://tiim.ch/blog/2022-09-27-sveltekit-ssr-with-urql, got %s", feed[0].uid)
		return
	}
}

func Test_wmSend_sendWebmentions(t *testing.T) {
	httpCalls := 0
	calledUrls := make([]string, 7)
	client := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method == "GET" {
			httpCalls += 1
			calledUrls = append(calledUrls, req.URL.String())
		}
		return &http.Response{
			StatusCode: 200,
		}
	})

	buf, err := ioutil.ReadFile("testdata/html/tiim.ch.rss-content.html")

	if err != nil {
		t.Errorf("unable to open testdata: %v", err)
	}

	wmsend := NewWmSend(&TestWMStore{}, client, "")
	now := time.Now()
	item := FeedItem{uid: "123", baseUrl: "https://tiim.ch/blog/2022-09-27-sveltekit-ssr-with-urql", updated: &now, content: string(buf)}
	err = wmsend.sendWebmentions(item)

	if err != nil {
		t.Fatal(err)
	}

	if httpCalls != 7 {
		t.Errorf("expected 7 http calls, got %d", httpCalls)
		return
	}
}

func Test_wmSend_sendWebmentions_preexisting_urls(t *testing.T) {
	httpCalls := 0
	calledUrls := make([]string, 7)
	client := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method == "GET" {
			httpCalls += 1
			calledUrls = append(calledUrls, req.URL.String())
		}
		return &http.Response{
			StatusCode: 200,
		}
	})

	buf, err := ioutil.ReadFile("testdata/html/tiim.ch.rss-content.html")

	if err != nil {
		t.Errorf("unable to open testdata: %v", err)
	}

	wmsend := NewWmSend(&TestWMStore{urls: []string{"https://example.com/1", "https://kit.svelte.dev/docs/load"}}, client, "")
	now := time.Now()
	item := FeedItem{uid: "123", baseUrl: "https://tiim.ch/blog/2022-09-27-sveltekit-ssr-with-urql", updated: &now, content: string(buf)}
	err = wmsend.sendWebmentions(item)

	if err != nil {
		t.Fatal(err)
	}

	if httpCalls != 8 {
		t.Errorf("expected 8 http calls, got %d", httpCalls)
		return
	}
}
