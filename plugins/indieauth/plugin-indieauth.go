package indieauth

import (
	"encoding/json"
	"fmt"
	"tiim/go-comment-api/config"
)

type indieAuthPlugin struct{}

type indieAuthPluginData struct {
	BaseUrl             string           `json:"baseUrl"`
	ProfileCanonicalUrl string           `json:"profileCanonicalUrl"`
	Password            string           `json:"password"`
	JWTSecret           string           `json:"jwtSecret"`
	StoreData           config.ModuleRaw `json:"store"`
}

func init() {
	config.RegisterPlugin(&indieAuthPlugin{})
}

func (p *indieAuthPlugin) Name() string {
	return "indieauth"
}

func (p *indieAuthPlugin) Load(data json.RawMessage, config config.GlobalConfig) (config.PluginInstance, error) {

	var d indieAuthPluginData
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}
	storeInt, err := config.Config.LoadModule(d.StoreData, nil)
	if err != nil {
		return nil, err
	}
	store, ok := storeInt.(Store)
	if !ok {
		return nil, fmt.Errorf("store module is not of type indieauth.Store: %T", storeInt)
	}

	return NewIndieAuthApiModule(d.BaseUrl, d.ProfileCanonicalUrl, d.Password, d.JWTSecret, store, *config.HttpClient), nil
}
