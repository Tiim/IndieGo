package webmentions_test

import (
	"os"
	"path"
	"testing"
	"tiim/go-comment-api/webmentions"

	"github.com/PuerkitoBio/goquery"
)

func parseHtml(name string) *goquery.Document {
	file, err := os.Open(path.Join("../test-data/html", name+".html"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		panic(err)
	}

	return doc
}

func TestEnrichContent(t *testing.T) {
	doc := parseHtml("webmention-rocks")
	wm := webmentions.Webmention{}
	enricher := webmentions.NewMicroformatEnricherChecker()
	err := enricher.CheckDocument(doc, &wm)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := "Test content."
	if wm.Content != expected {
		t.Errorf("unexpected content: \n%s\n%s", wm.Content, expected)
	}
}

func TestEnrichContentSummary(t *testing.T) {
	doc := parseHtml("webmention-rocks-summary")
	wm := webmentions.Webmention{}
	enricher := webmentions.NewMicroformatEnricherChecker()
	err := enricher.CheckDocument(doc, &wm)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := "Make it big!!!"
	if wm.Content != expected {
		t.Errorf("unexpected content: \n%s\n%s", wm.Content, expected)
	}
}
