package mfobjects

import (
	"time"

	"willnorris.com/go/microformats"
)

type MF2HEntry struct {
	Name        string    `yaml:"name,omitempty"`
	Summary     string    `yaml:"summary,omitempty"`
	Content     string    `yaml:"-"`
	Published   time.Time `yaml:"date,omitempty"`
	Updated     time.Time `yaml:"modified,omitempty"`
	Author      MF2HCard  `yaml:"author,omitempty"`
	Category    []string  `yaml:"content_tags,omitempty"`
	Url         string    `yaml:"-"`
	Photos      MF2Photos `yaml:"photos,omitempty"`
	InReplyTo   MF2HCite  `yaml:"in_reply_to,omitempty"`
	RSVP        string    `yaml:"rsvp,omitempty"`
	LikeOf      MF2HCite  `yaml:"like_of,omitempty"`
	RepostOf    MF2HCite  `yaml:"repost_of,omitempty"`
	Syndication []string  `yaml:"syndication,omitempty"`
}

type MF2HCite struct {
	Name        string    `yaml:"name,omitempty"`
	Published   time.Time `yaml:"published,omitempty"`
	Author      MF2HCard  `yaml:"author,omitempty"`
	Url         string    `yaml:"url,omitempty"`
	Publication string    `yaml:"publication,omitempty"`
	Accessed    time.Time `yaml:"accessed,omitempty"`
	Content     string    `yaml:"content,omitempty"`
	Summary     string    `yaml:"summary,omitempty"`
}

type MF2HCard struct {
	Name string `yaml:"name,omitempty"`
}

type MF2HApp struct {
	Url          string   `yaml:"url,omitempty"`
	Name         string   `yaml:"name,omitempty"`
	Logo         string   `yaml:"logo,omitempty"`
	Summary      string   `yaml:"summary,omitempty"`
	Author       MF2HCard `yaml:"author,omitempty"`
	RedirectUris []string `yaml:"redirect_uris,omitempty"`
}

type MF2Photo struct {
	Url string `yaml:"url,omitempty"`
	Alt string `yaml:"alt,omitempty"`
}

type MF2Photos []MF2Photo

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
	if len(h.Syndication) > 0 {
		mf.Properties["syndication"] = []interface{}{}
		for _, syndication := range h.Syndication {
			mf.Properties["syndication"] = append(mf.Properties["syndication"], syndication)
		}
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
