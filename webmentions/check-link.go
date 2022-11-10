package webmentions

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type linkToTargetChecker struct{}

func (c *linkToTargetChecker) CheckMention(w *Webmention) error {
	return nil
}
func (c *linkToTargetChecker) CheckDocument(gq *goquery.Document, w *Webmention) error {
	links := gq.Find("a")
	valid := false
	links.Each(func(i int, s *goquery.Selection) {
		href, exits := s.Attr("href")
		if exits && href == w.Target {
			valid = true
		}
	})
	if valid {
		return nil
	} else {
		return fmt.Errorf("no link to target")
	}
}

func NewLinkToTargetChecker() *linkToTargetChecker {
	return &linkToTargetChecker{}
}
