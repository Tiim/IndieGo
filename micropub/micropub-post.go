package micropub

import (
	"tiim/go-comment-api/mfobjects"
	"time"

	"willnorris.com/go/microformats"
)

type MicropubPostRaw struct {
	Action      string                   `json:"action"`
	Url         string                   `json:"url"`
	PostTye     []string                 `json:"type"`
	Properties  map[string][]interface{} `json:"properties"`
	Add         map[string][]interface{} `json:"add"`
	Replace     map[string][]interface{} `json:"replace"`
	Delete      interface{}              `json:"delete"`
	AccessToken string                   `json:"-"`
}

type MicropubPost struct {
	Entry mfobjects.MF2HEntry
}

func ParseMicropubPost(data MicropubPostRaw) MicropubPost {
	mf := microformats.Data{
		Items: []*microformats.Microformat{
			{
				Type:       data.PostTye,
				Properties: data.Properties,
			},
		},
	}
	post := MicropubPost{
		Entry: mfobjects.GetHEntry(&mf),
	}
	if post.Entry.Published.IsZero() {
		post.Entry.Published = time.Now().UTC()
	}
	return post
}

func (post *MicropubPost) ToMarkdown() string {
	return post.Entry.ToMarkdown()
}

func PostFromMarkdown(markdown string) MicropubPost {
	return MicropubPost{
		Entry: mfobjects.EntryFromMarkdonw(markdown),
	}
}
