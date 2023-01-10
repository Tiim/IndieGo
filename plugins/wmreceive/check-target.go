package wmrecv

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type targetChecker struct {
	domain []string
}

func (c *targetChecker) CheckMention(w *Webmention) error {
	if !strInSlice(w.TargetUrl().Hostname(), c.domain) {
		return fmt.Errorf("target domain is not in allow list")
	}
	return nil
}

func (c *targetChecker) CheckDocument(gq *goquery.Document, w *Webmention) error {
	return nil
}

func newTargetChecker(domain ...string) *targetChecker {
	return &targetChecker{domain}
}

func strInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}
