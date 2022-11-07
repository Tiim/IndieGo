package webmentions

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Checker interface {
	CheckMention(w *Webmention) error
	CheckDocument(gq *goquery.Document, w *Webmention) error
}

type WebmentionChecker struct {
	checkers []Checker
}

func NewWebmentionChecker(checkers []Checker) *WebmentionChecker {
	return &WebmentionChecker{checkers: checkers}
}

func (c *WebmentionChecker) CheckWebmentionValid(w *Webmention) error {
	errors := make([]error, 0)

	for _, checker := range c.checkers {
		if err := checker.CheckMention(w); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("webmention failed checks: %v", errors)
	}

	res, err := http.Get(w.Source)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	html, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	for _, checker := range c.checkers {
		err := checker.CheckDocument(html, w)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("webmention failed checks: %v", errors)
	}

	return nil
}

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

type domainChecker struct {
	store *webmentionsStore
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

func NewDomainChecker(store *webmentionsStore) *domainChecker {
	return &domainChecker{store}
}
