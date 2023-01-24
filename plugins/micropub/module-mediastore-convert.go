package micropub

import (
	"fmt"
	"log"
	"tiim/go-comment-api/config"
)

type MediastoreConvertModule struct {
	FormatMap  map[string]string `json:"format_map"`
	MediaStore config.ModuleRaw  `json:"media_store" config:"micropub.media-store"`
}

func init() {
	config.RegisterModule(&MediastoreConvertModule{})
}

func (m *MediastoreConvertModule) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "micropub.media-store.convert",
		New:  func() config.Module { return new(MediastoreConvertModule) },
		Docs: config.ConfigDocs{
			DocString: `Media store moduel that converts files from one media file format into another before passing it to another media store module.`,
			Fields: map[string]string{
				"FormatMap": `JSON map of source -> destination mime types. 
					Example <code>{\"image/jpeg\": \"image/webp\"}</code>
					The special key "*" can be used to match as a fallback for all mime types. The special value '-' means no conversion.
					The default value is {"*": "-"} which does not convert anything.
					Currently the following mime types are supported: 
					<ul><li>image/jpeg</li><li>image/png</li><li>image/webp</li><li>image/gif</li></ul>
					`,
				"MediaStore": `The media store module to pass the converted file to.`,
			},
		},
	}
}

func (m *MediastoreConvertModule) Load(config config.GlobalConfig, args interface{}, logger *log.Logger) (config.ModuleInstance, error) {
	ms, err := config.Config.LoadModule(m, "MediaStore", args)
	if err != nil {
		return nil, fmt.Errorf("failed to load media store module: %w", err)
	}
	return &convertMediaStore{
		childMediaStore: ms.(mediaStore),
		convertMap:      m.FormatMap,
		logger:          logger,
	}, nil
}
