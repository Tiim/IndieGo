package micropub

import (
	"encoding/json"
	"strings"
	"tiim/go-comment-api/config"
)

type MediastoreStorjModule struct{}
type MediastoreStorjModuleData struct {
	AccessGrant string `json:"access_grant"`
	BucketName  string `json:"bucket_name"`
	Prefix      string `json:"prefix"`

	// The format of the url to the media file:
	// {name} will be replaced with the name of the file,
	// {prefix} will be replaced with the prefix,
	// {bucket} will be replaced with the bucket name.
	UrlFormat string `json:"url_format"`
}

func init() {
	config.RegisterModule(&MediastoreStorjModule{})
}

func (p *MediastoreStorjModule) Name() string {
	return "micropub-mediastore-storj"
}

func (p *MediastoreStorjModule) Load(data json.RawMessage, config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	d := MediastoreStorjModuleData{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}

	return newStorjMediaStore(
		d.AccessGrant,
		d.BucketName,
		d.Prefix,
		func(name, contentType, prefix, bucket string) string {
			url := d.UrlFormat
			url = strings.Replace(url, "{name}", name, -1)
			url = strings.Replace(url, "{prefix}", prefix, -1)
			url = strings.Replace(url, "{bucket}", bucket, -1)
			return url
		},
	), nil
}
