package micropub

import (
	"strings"
	"tiim/go-comment-api/config"
)

type MediastoreStorjModule struct {
	// The storj access grant.
	AccessGrant string `json:"access_grant"`
	// The name of the storj bucket.
	BucketName string `json:"bucket_name"`
	// Can be a custom prefix or an empty string for no prefix.
	// Setting a prefix allows multiple uses of the same bucket.
	Prefix string `json:"prefix"`

	// The format of the url to the media file:
	// {name} will be replaced with the name of the file,
	// {prefix} will be replaced with the prefix,
	// {bucket} will be replaced with the bucket name.
	UrlFormat string `json:"url_format"`
}

func init() {
	config.RegisterModule(&MediastoreStorjModule{})
}

func (m *MediastoreStorjModule) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "micropub.media-store.storj",
		New:  func() config.Module { return new(MediastoreStorjModule) },
	}
}

func (m *MediastoreStorjModule) Load(config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	return newStorjMediaStore(
		m.AccessGrant,
		m.BucketName,
		m.Prefix,
		func(name, contentType, prefix, bucket string) string {
			url := m.UrlFormat
			url = strings.Replace(url, "{name}", name, -1)
			url = strings.Replace(url, "{prefix}", prefix, -1)
			url = strings.Replace(url, "{bucket}", bucket, -1)
			return url
		},
	), nil
}
