package webmentions

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type targetChecker struct {
	domain []string
}

func (c *targetChecker) CheckMention(w *Webmention) error {
	if !strInArray(w.TargetUrl().Hostname(), c.domain) {
		return fmt.Errorf("target domain is not in allow list")
	}
	return nil
}

func (c *targetChecker) CheckDocument(gq *goquery.Document, w *Webmention) error {
	return nil
}

func NewTargetChecker(domain ...string) *targetChecker {
	return &targetChecker{domain}
}
