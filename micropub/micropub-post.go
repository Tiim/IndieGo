package micropub

import (
	"tiim/go-comment-api/mfobjects"

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
	return MicropubPost{
		Entry: mfobjects.GetHEntry(&mf),
	}
}
