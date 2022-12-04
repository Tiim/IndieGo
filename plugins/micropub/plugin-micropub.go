package micropub

import (
	"encoding/json"
	"fmt"
	"log"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/plugins/indieauth"
)

type micropubPlugin struct{}
type micropubPluginData struct {
	StoreData      config.ModuleRaw `json:"store"`
	MediaStoreData config.ModuleRaw `json:"media_store"`
}

func init() {
	config.RegisterPlugin(&micropubPlugin{})
}

func (p *micropubPlugin) Name() string {
	return "micropub"
}

func (p *micropubPlugin) Load(data json.RawMessage, config config.GlobalConfig) (config.PluginInstance, error) {
	d := micropubPluginData{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}

	storeInt, err := config.Config.LoadModule(d.StoreData, nil)
	if err != nil {
		return nil, err
	}
	store, ok := storeInt.(micropubStore)
	if !ok {
		return nil, fmt.Errorf("store module is not of type micropub.micropubStore: %T", storeInt)
	}

	mstoreInt, err := config.Config.LoadModule(d.MediaStoreData, nil)
	if err != nil {
		return nil, err
	}
	mstore, ok := mstoreInt.(mediaStore)
	if !ok {
		return nil, fmt.Errorf("media store module is not of type micropub.mediaStore: %T", mstoreInt)
	}

	indieAuthPlugin, err := config.GetPlugin("indieauth")
	if err != nil {
		log.Println("The micropub plugin requires the indieauth plugin to be loaded. If you want to verify a token from another source, please open an issue on github.")
		return nil, err
	}
	indieAuth, ok := indieAuthPlugin.(*indieauth.IndieAuthApiModule)

	return newMicropubApiModule(store, mstore, indieAuth.VerifyToken), nil
}
