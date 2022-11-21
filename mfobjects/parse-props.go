package mfobjects

import (
	"fmt"
	"strings"
	"time"

	"willnorris.com/go/microformats"
)

type MF2HEntry struct {
	Name      string
	Summary   string
	Content   string
	Published time.Time
	Updated   time.Time
	Author    MF2HCard
	Category  []string
	Url       string
	Photos    MF2Photos
	InReplyTo MF2HCite
	RSVP      string
	LikeOf    MF2HCite
	RepostOf  MF2HCite
}

type MF2HCite struct {
	Name        string
	Published   time.Time
	Author      MF2HCard
	Url         string
	Publication string
	Accessed    time.Time
	Content     string
	Summary     string
}

type MF2HCard struct {
	Name string
}

type MF2HApp struct {
	Url          string
	Name         string
	Logo         string
	Summary      string
	Author       MF2HCard
	RedirectUris []string
}

type MF2Photo struct {
	Url string
	Alt string
}

type MF2Photos []MF2Photo

func GetHApp(data *microformats.Data) MF2HApp {
	item := getEntryWithType(data, "h-app", "h-x-app")
	if item == nil {
		return MF2HApp{}
	}
	return MF2HApp{
		Url:          GetStringProp("url", item),
		Name:         GetStringProp("name", item),
		Logo:         GetStringProp("logo", item),
		Summary:      GetStringProp("summary", item),
		Author:       GetHCard("author", item),
		RedirectUris: GetStringPropSlice("redirect-uri", item),
	}
}

func GetHEntry(data *microformats.Data) MF2HEntry {
	item := getEntryWithType(data, "h-entry")
	if item == nil {
		return MF2HEntry{}
	}

	return MF2HEntry{
		Name:      GetStringProp("name", item),
		Summary:   GetStringProp("summary", item),
		Content:   GetStringProp("content", item),
		Published: GetTimeProp("published", item),
		Updated:   GetTimeProp("updated", item),
		Author:    GetHCard("author", item),
		Category:  GetStringPropSlice("category", item),
		Url:       GetStringProp("url", item),
		Photos:    GetPhotos("photo", item),
		InReplyTo: GetHCite("in-reply-to", item),
		LikeOf:    GetHCite("like-of", item),
		RepostOf:  GetHCite("repost-of", item),
	}
}

func GetHCard(name string, item *microformats.Microformat) MF2HCard {
	author := item.Properties[name]
	if len(author) > 0 {
		if authorStr, ok := author[0].(string); ok {
			return MF2HCard{Name: authorStr}
		}
		authorMf, ok := author[0].(*microformats.Microformat)
		if ok {
			return MF2HCard{Name: authorMf.Value}
		}
	}
	return MF2HCard{}
}

func GetHCite(name string, item *microformats.Microformat) MF2HCite {
	cite, ok := item.Properties[name]
	if !ok {
		return MF2HCite{}
	}
	if len(cite) == 0 {
		return MF2HCite{}
	}
	if citeString, ok := cite[0].(string); ok {
		return MF2HCite{Url: citeString}
	}
	citeMf, ok := cite[0].(*microformats.Microformat)
	if ok {
		return MF2HCite{
			Name:        GetStringProp("name", citeMf),
			Published:   GetTimeProp("published", citeMf),
			Author:      GetHCard("author", citeMf),
			Url:         GetStringProp("url", citeMf),
			Publication: GetStringProp("publication", citeMf),
			Accessed:    GetTimeProp("accessed", citeMf),
			Content:     GetStringProp("content", citeMf),
			Summary:     GetStringProp("summary", citeMf),
		}
	}
	return MF2HCite{}
}

func GetPhotos(name string, item *microformats.Microformat) MF2Photos {
	propValue, ok := item.Properties[name]
	if !ok || len(propValue) == 0 {
		return []MF2Photo{}
	}
	slice := make([]MF2Photo, 0, len(propValue))
	for _, val := range propValue {
		if value, ok := val.(string); ok {
			slice = append(slice, MF2Photo{Url: value})
		} else if value, ok := val.(map[string]interface{}); ok {
			var url string
			var alt string
			if value["value"] != nil && value["value"].(string) != "" {
				url = value["value"].(string)
			}
			if value["alt"] != nil && value["alt"].(string) != "" {
				alt = value["alt"].(string)
			}
			slice = append(slice, MF2Photo{Url: url, Alt: alt})
		}
	}
	return slice
}

func getEntryWithType(data *microformats.Data, types ...string) *microformats.Microformat {
	for _, item := range data.Items {
		for _, itemType := range item.Type {
			for _, searchType := range types {
				if itemType == searchType {
					return item
				}
			}
		}
	}
	return nil
}

func GetStringProp(name string, item *microformats.Microformat) string {
	propValue, ok := item.Properties[name]
	if !ok {
		return ""
	}
	if len(propValue) == 0 {
		return ""
	}
	if value, ok := propValue[0].(string); ok {
		return trimLines(value)
	}
	if value, ok := propValue[0].(map[string]interface{}); ok {
		if value["value"] != nil && value["value"].(string) != "" {
			return trimLines(value["value"].(string))
		} else if value["html"] != nil && value["html"].(string) != "" {
			return trimLines(value["html"].(string))
		}
	}
	fmt.Printf("Did not find string prop %s: %v (%T)", name, propValue[0], propValue[0])
	return ""
}

func GetStringPropSlice(name string, item *microformats.Microformat) []string {
	propValue, ok := item.Properties[name]
	if !ok {
		return []string{}
	}
	if len(propValue) == 0 {
		return []string{}
	}
	slice := make([]string, 0, len(propValue))
	for _, value := range propValue {
		if value, ok := value.(string); ok {
			slice = append(slice, trimLines(value))
		}
	}
	return slice
}

var timeFormats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05-07:00",
	"2006-01-02T15:04:05-0700",
}

func GetTimeProp(name string, item *microformats.Microformat) time.Time {
	propValue, ok := item.Properties[name]
	if !ok {
		return time.Time{}
	}
	if len(propValue) == 0 {
		return time.Time{}
	}
	if value, ok := propValue[0].(string); ok {
		var parsed time.Time
		for _, format := range timeFormats {
			var err error
			parsed, err = time.Parse(format, value)
			if err == nil {
				break
			}
		}
		if parsed.IsZero() {
			return time.Time{}
		}
		parsed = parsed.UTC()
		return parsed
	}
	fmt.Printf("Did not find time prop %s: %v (%T)", name, propValue[0], propValue[0])
	return time.Time{}
}

func (h *MF2HEntry) GetShortContent(maxLength, maxNewlines int) string {
	if h.Content != "" && !isContentTooLong(h.Content, maxLength, maxNewlines) {
		return h.Content
	} else if h.Summary != "" && !isContentTooLong(h.Summary, maxLength, maxNewlines) {
		return h.Summary
	} else if h.Summary != "" {
		return truncateContent(h.Summary, maxLength, maxNewlines)
	} else if h.Content != "" {
		return truncateContent(h.Content, maxLength, maxNewlines)
	} else if h.Name != "" {
		return truncateContent(h.Name, maxLength, maxNewlines)
	}
	return ""
}

func truncateContent(content string, maxLength, maxNewlines int) string {
	truncated := false
	if strings.Count(content, "\n") > maxNewlines {
		content = strings.Join(strings.Split(content, "\n")[0:maxNewlines], "\n")
		truncated = true
	}
	if len(content) > maxLength {
		content = content[0:maxLength]
		truncated = true
	}
	if truncated {
		content = content + "..."
	}
	return content
}

func isContentTooLong(content string, maxLength, maxNewlines int) bool {
	return len(content) > maxLength || strings.Count(content, "\n") > maxNewlines
}

func trimLines(s string) string {
	s = strings.TrimSpace(s)
	lines := strings.Split(s, "\n")
	newLines := make([]string, 0)
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if trim != "" {
			newLines = append(newLines, trim)
		}
	}
	return strings.Join(newLines, "\n")
}

func (h *MF2HEntry) ToMicroformat() *microformats.Microformat {
	mf := &microformats.Microformat{
		Type:       []string{"h-entry"},
		Properties: map[string][]interface{}{},
	}
	if h.Name != "" {
		mf.Properties["name"] = []interface{}{h.Name}
	}
	if h.Summary != "" {
		mf.Properties["summary"] = []interface{}{h.Summary}
	}
	if h.Content != "" {
		mf.Properties["content"] = []interface{}{h.Content}
	}
	if !h.Published.IsZero() {
		mf.Properties["published"] = []interface{}{h.Published.Format(time.RFC3339)}
	}
	if !h.Updated.IsZero() {
		mf.Properties["updated"] = []interface{}{h.Updated.Format(time.RFC3339)}
	}
	if h.Author.Name != "" {
		mf.Properties["author"] = []interface{}{h.Author.ToMicroformat()}
	}
	if len(h.Category) > 0 {
		mf.Properties["category"] = []interface{}{}
		for _, category := range h.Category {
			mf.Properties["category"] = append(mf.Properties["category"], category)
		}
	}
	if h.Url != "" {
		mf.Properties["url"] = []interface{}{h.Url}
	}
	if len(h.Photos) > 0 {
		mf.Properties["photo"] = []interface{}{h.Photos.ToMicroformat()}
	}
	if h.InReplyTo.Url != "" {
		mf.Properties["in-reply-to"] = []interface{}{h.InReplyTo.ToMicroformat()}
	}
	if h.RSVP != "" {
		mf.Properties["rsvp"] = []interface{}{h.RSVP}
	}
	if h.LikeOf.Url != "" {
		mf.Properties["like-of"] = []interface{}{h.LikeOf.ToMicroformat()}
	}
	if h.RepostOf.Url != "" {
		mf.Properties["repost-of"] = []interface{}{h.RepostOf.ToMicroformat()}
	}

	return mf
}

func (h *MF2HCard) ToMicroformat() *microformats.Microformat {
	mf := &microformats.Microformat{
		Type:       []string{"h-card"},
		Properties: map[string][]interface{}{},
	}
	if h.Name != "" {
		mf.Properties["name"] = []interface{}{h.Name}
	}
	return mf
}

func (h *MF2HCite) ToMicroformat() *microformats.Microformat {
	mf := &microformats.Microformat{
		Type:       []string{"h-cite"},
		Properties: map[string][]interface{}{},
	}
	if h.Name != "" {
		mf.Properties["name"] = []interface{}{h.Name}
	}
	if !h.Published.IsZero() {
		mf.Properties["published"] = []interface{}{h.Published.Format(time.RFC3339)}
	}
	if h.Author.Name != "" {
		mf.Properties["author"] = []interface{}{h.Author.ToMicroformat()}
	}
	if h.Url != "" {
		mf.Properties["url"] = []interface{}{h.Url}
	}
	if h.Publication != "" {
		mf.Properties["publication"] = []interface{}{h.Publication}
	}
	if !h.Accessed.IsZero() {
		mf.Properties["accessed"] = []interface{}{h.Accessed.Format(time.RFC3339)}
	}
	if h.Content != "" {
		mf.Properties["content"] = []interface{}{h.Content}
	}
	if h.Summary != "" {
		mf.Properties["summary"] = []interface{}{h.Summary}
	}
	return mf
}

func (p *MF2Photos) ToMicroformat() []map[string]interface{} {
	slice := make([]map[string]interface{}, 0)
	for _, photo := range *p {
		slice = append(slice, map[string]interface{}{"url": photo.Url, "alt": photo.Alt})
	}
	return slice
}
