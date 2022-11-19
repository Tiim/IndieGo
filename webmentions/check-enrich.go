package webmentions

import (
	"tiim/go-comment-api/mfobjects"

	"github.com/PuerkitoBio/goquery"
	"willnorris.com/go/microformats"
)

type microformatEnricherChecker struct{}

func (c *microformatEnricherChecker) CheckMention(w *Webmention) error {
	return nil
}

func (c *microformatEnricherChecker) CheckDocument(gq *goquery.Document, w *Webmention) error {
	data := microformats.ParseNode(gq.Nodes[0], w.SourceUrl())
	hentry := mfobjects.GetHEntry(data)

	w.AuthorName = hentry.Author.Name
	w.Content = hentry.GetShortContent(500, 4)

	return nil
}

func NewMicroformatEnricherChecker() *microformatEnricherChecker {
	return &microformatEnricherChecker{}
}
