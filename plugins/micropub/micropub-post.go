package micropub

import (
	"fmt"
	"io"
	"log"
	"strings"
	"tiim/go-comment-api/lib/mfobjects"
	"time"

	"gopkg.in/yaml.v3"
	"willnorris.com/go/microformats"
)

type MicropubFile struct {
	ContentType string
	Reader      io.ReadCloser
	Name        string
	Size        int64
}

type MicropubPostRaw struct {
	Action     string                   `json:"action"`
	Url        string                   `json:"url"`
	PostTye    []string                 `json:"type"`
	Properties map[string][]interface{} `json:"properties"`
	Add        map[string][]interface{} `json:"add"`
	Replace    map[string][]interface{} `json:"replace"`
	Delete     interface{}              `json:"delete"`
	Files      []MicropubFile           `json:"-"`
}

type MicropubPost struct {
	Entry   mfobjects.MF2HEntry `yaml:",inline"`
	RawData *microformats.Data  `yaml:"raw_data,flow,omitempty"`
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
		Entry:   mfobjects.GetHEntry(&mf),
		RawData: &mf,
	}
	if post.Entry.Published.IsZero() {
		post.Entry.Published = time.Now().UTC()
	}
	return post
}

func (post *MicropubPost) ToMarkdown() string {
	frontmatter, err := yaml.Marshal(post)
	if err != nil {
		log.Println("Error marshalling raw microformat data to JSON: ", err)
		frontmatter = []byte("{}")
	}
	return fmt.Sprintf("---\n%s---\n\n%s", frontmatter, post.Entry.Content)
}

func PostFromMarkdown(markdown string) MicropubPost {
	splits := strings.Split(markdown, "---")
	if len(splits) < 3 {
		log.Println("Markdown does not contain frontmatter")
		return MicropubPost{
			Entry: mfobjects.MF2HEntry{
				Content: markdown,
			},
			RawData: &microformats.Data{},
		}
	}
	frontmatterStr := splits[1]
	var frontmatter MicropubPost
	err := yaml.Unmarshal([]byte(frontmatterStr), &frontmatter)
	if err != nil {
		log.Println("Could not parse frontmatter", err)
		return MicropubPost{
			Entry: mfobjects.MF2HEntry{
				Content: markdown,
			},
			RawData: &microformats.Data{},
		}
	}

	return frontmatter
}
