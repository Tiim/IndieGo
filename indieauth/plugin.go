package indieauth

import (
	"encoding/json"
	"fmt"
	"tiim/go-comment-api/plugin"
)

type indieAuthPlugin struct{}

type indieAuthPluginData struct {
	BaseUrl             string           `json:"baseUrl"`
	ProfileCanonicalUrl string           `json:"profileCanonicalUrl"`
	Password            string           `json:"password"`
	JWTSecret           string           `json:"jwtSecret"`
	StoreData           plugin.ModuleRaw `json:"store"`
}

func init() {
	plugin.RegisterPlugin(&indieAuthPlugin{})
}

func (p *indieAuthPlugin) Name() string {
	return "indieauth"
}

func (p *indieAuthPlugin) Load(data json.RawMessage, config plugin.GlobalConfig) (plugin.PluginInstance, error) {

	var d indieAuthPluginData
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}
	storeInt, err := config.Config.LoadModule(d.StoreData)
	if err != nil {
		return nil, err
	}
	store, ok := storeInt.(Store)
	if !ok {
		return nil, fmt.Errorf("store module is not of type indieauth.Store: %T", storeInt)
	}

	return NewIndieAuthApiModule(d.BaseUrl, d.ProfileCanonicalUrl, d.Password, d.JWTSecret, store, *config.HttpClient), nil
}
