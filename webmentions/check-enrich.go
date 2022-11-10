package webmentions

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"willnorris.com/go/microformats"
)

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
	eContent := trimLines(getMF2EContent(hentry))
	pSummary := trimLines(getMF2PSummary(hentry))
	pName := trimLines(getMF2PName(hentry))
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

func trimLines(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n\n", "\n")
	lines := strings.Split(s, "\n")
	newLines := make([]string, len(lines))
	for i, line := range lines {
		trim := strings.TrimSpace(line)
		if trim != "" {
			newLines[i] = trim
		}
	}
	return strings.Join(newLines, "\n")
}

func getMF2EContent(hentry *microformats.Microformat) string {
	content := hentry.Properties["content"]
	if len(content) > 0 {
		contentMf, ok := content[0].(map[string]string)
		if ok {
			return contentMf["value"]
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
