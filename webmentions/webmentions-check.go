package webmentions

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"willnorris.com/go/microformats"
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

type microformatEnricherChecker struct{}

func (c *microformatEnricherChecker) CheckMention(w *Webmention) error {
	return nil
}

func (c *microformatEnricherChecker) CheckDocument(gq *goquery.Document, w *Webmention) error {
	data := microformats.ParseNode(gq.Nodes[0], w.SourceUrl())
	for _, item := range data.Items {
		if strInArray("h-entry", item.Type) {
			extractAuthor(item, w)
			extractContent(item, w)
		}
	}
	return nil
}

func strInArray(str string, arr []string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func extractAuthor(hentry *microformats.Microformat, w *Webmention) {
	author := hentry.Properties["author"]
	if len(author) > 0 {
		authorMf, ok := author[0].(*microformats.Microformat)
		if ok {
			w.AuthorName = authorMf.Value
		}
	}
}

func extractContent(hentry *microformats.Microformat, w *Webmention) {
	reply := hentry.Properties["in-reply-to"]
	isReply := false
	for _, replyTo := range reply {
		replyMf, ok := replyTo.(*microformats.Microformat)
		if ok {
			if replyMf.Value == w.Target {
				isReply = true
			}
		}
	}
	if !isReply {
		w.Content = ""
		return
	}
	eContent := getMF2EContent(hentry)
	pSummary := getMF2PSummary(hentry)
	pName := getMF2PName(hentry)
	w.Content = truncateContentAndSummary(eContent, pSummary, pName)
}

func truncateContentAndSummary(content, summary, name string) string {
	if content != "" && !isContentTooLong(content) {
		return content
	} else if summary != "" && !isContentTooLong(summary) {
		return summary
	} else if summary != "" {
		return truncateContent(summary)
	} else if content != "" {
		return truncateContent(content)
	} else if name != "" {
		return truncateContent(name)
	}
	return ""
}

func isContentTooLong(content string) bool {
	return len(content) > 500 || strings.Count(content, "\n") > 4
}
func truncateContent(content string) string {
	truncated := false
	if strings.Count(content, "\n") > 4 {
		content = strings.Join(strings.Split(content, "\n")[0:4], "\n")
		truncated = true
	}
	if len(content) > 500 {
		content = content[0:500]
		truncated = true
	}
	if truncated {
		content = content + "..."
	}

	return content
}

func getMF2EContent(hentry *microformats.Microformat) string {
	content := hentry.Properties["content"]
	if len(content) > 0 {
		contentMf, ok := content[0].(*microformats.Microformat)
		if ok {
			return contentMf.Value
		}
	}
	return ""
}

func getMF2PSummary(hentry *microformats.Microformat) string {
	summary := hentry.Properties["summary"]
	if len(summary) > 0 {
		summary, ok := summary[0].(string)
		if ok {
			return summary
		}
	}
	return ""
}

func getMF2PName(hentry *microformats.Microformat) string {
	name := hentry.Properties["name"]
	if len(name) > 0 {
		name, ok := name[0].(string)
		if ok {
			return name
		}
	}
	return ""
}

func NewMicroformatEnricherChecker() *microformatEnricherChecker {
	return &microformatEnricherChecker{}
}
