package micropub

import (
	"fmt"
	"log"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/plugins/indieauth"
)

type micropubPlugin struct {
	StoreData      config.ModuleRaw `json:"store" config:"micropub.store"`
	MediaStoreData config.ModuleRaw `json:"media_store" config:"micropub.media-store"`
}

func init() {
	config.RegisterModule(&micropubPlugin{})
}

func (p *micropubPlugin) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "micropub",
		New:  func() config.Module { return new(micropubPlugin) },
		Docs: config.ConfigDocs{
			DocString: `Micropub module. This module enables the micropub endpoint.`,
			Fields: map[string]string{
				"StoreData":      "The store module to use for storing micropub data.",
				"MediaStoreData": "The media store module to use for storing media.",
			},
		},
	}
}

func (p *micropubPlugin) Load(config config.GlobalConfig, _ interface{}) (config.ModuleInstance, error) {

	storeInt, err := config.Config.LoadModule(p, "StoreData", nil)
	if err != nil {
		return nil, err
	}
	store, ok := storeInt.(micropubStore)
	if !ok {
		return nil, fmt.Errorf("store module is not of type micropub.micropubStore: %T", storeInt)
	}

	mstoreInt, err := config.Config.LoadModule(p, "MediaStoreData", nil)
	if err != nil {
		return nil, err
	}
	mstore, ok := mstoreInt.(mediaStore)
	if !ok {
		return nil, fmt.Errorf("media store module is not of type micropub.mediaStore: %T", mstoreInt)
	}

	indieAuthPlugin, err := config.GetModule("indieauth")
	if err != nil {
		log.Println("The micropub plugin requires the indieauth plugin to be loaded. If you want to verify a token from another source, please open an issue on github.")
		return nil, err
	}
	indieAuth, ok := indieAuthPlugin.(*indieauth.IndieAuthApiModule)
	if !ok {
		return nil, fmt.Errorf("indieauth plugin is not of type indieauth.IndieAuthApiModule: %T", indieAuthPlugin)
	}

	return newMicropubApiModule(store, mstore, indieAuth.VerifyToken), nil
}
