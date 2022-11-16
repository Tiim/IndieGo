package webmentions

import (
	"tiim/go-comment-api/microformatsextract"

	"github.com/PuerkitoBio/goquery"
	"willnorris.com/go/microformats"
)

type microformatEnricherChecker struct{}

func (c *microformatEnricherChecker) CheckMention(w *Webmention) error {
	return nil
}

func (c *microformatEnricherChecker) CheckDocument(gq *goquery.Document, w *Webmention) error {
	data := microformats.ParseNode(gq.Nodes[0], w.SourceUrl())
	hentry := microformatsextract.GetHEntry(data)

	if hentry != nil {
		if hentry.Author != nil {
			w.AuthorName = hentry.Author.Name
		}
		w.Content = hentry.GetShortContent(500, 4)
	}

	return nil
}

func NewMicroformatEnricherChecker() *microformatEnricherChecker {
	return &microformatEnricherChecker{}
}
