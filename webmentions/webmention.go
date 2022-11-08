package webmentions

import (
	"fmt"
	"net/url"
	"strings"
	"tiim/go-comment-api/model"
	"time"

	"github.com/google/uuid"
)

type Webmention struct {
	Id        string
	Source    string
	Target    string
	TsCreated time.Time
	TsUpdated time.Time
}

func NewWebmention(source, target string) (*Webmention, error) {

	sourceUrl, err := url.ParseRequestURI(source)

	if err != nil || !strings.HasPrefix(sourceUrl.Scheme, "http") {
		return nil, fmt.Errorf("invalid source url: %w", err)
	}

	targetUrl, err := url.ParseRequestURI(target)

	if err != nil || !strings.HasPrefix(targetUrl.Scheme, "http") {
		return nil, fmt.Errorf("invalid target url: %w", err)
	}

	if *sourceUrl == *targetUrl {
		return nil, fmt.Errorf("source and target are the same")
	}

	return &Webmention{
		Id:        uuid.New().String(),
		Source:    source,
		Target:    target,
		TsCreated: time.Now(),
		TsUpdated: time.Now(),
	}, nil
}

func (w *Webmention) ToGenericComment() model.GenericComment {
	c := model.GenericComment{
		Id:        w.Id,
		Type:      "webmention",
		Timestamp: w.TsCreated.Format(time.RFC3339),
		Page:      w.Page(),
		Content:   w.Source,
	}
	return c
}

func (w *Webmention) SourceUrl() *url.URL {
	u, _ := url.Parse(w.Source)
	return u
}

func (w *Webmention) Page() string {
	u, _ := url.Parse(w.Target)
	page := u.Path
	if page[0] == '/' {
		page = page[1:]
	}
	return page
}
