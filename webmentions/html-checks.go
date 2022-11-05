package webmentions

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Checker func(gq *goquery.Document, w *Webmention) error

var (
	checkers = []Checker{
		linkToTarget,
	}
)

func checkWebmentionValid(w *Webmention) error {
	res, err := http.Get(w.Source)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	html, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	errors := make([]error, 0)
	for _, checker := range checkers {
		err := checker(html, w)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("webmention failed checks: %v", errors)
	}

	return nil
}

func linkToTarget(gq *goquery.Document, w *Webmention) error {
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
