package webmentions

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type domainChecker struct {
	store webmentionsStore
}

func (c *domainChecker) CheckMention(w *Webmention) error {
	denylist, err := c.store.GetDomainDenyList()
	if err != nil {
		return err
	}

	for _, domain := range denylist {
		if w.SourceUrl().Hostname() == domain {
			return fmt.Errorf("source domain is in deny list")
		}
	}
	return nil
}
func (c *domainChecker) CheckDocument(gq *goquery.Document, w *Webmention) error {
	return nil
}

func NewDomainChecker(store webmentionsStore) *domainChecker {
	return &domainChecker{store}
}
