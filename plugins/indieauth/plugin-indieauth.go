package indieauth

import (
	"fmt"
	"strings"
	"tiim/go-comment-api/config"
)

type indieAuthPlugin struct {
	// The url indiego is running on. For example https://indiego.example.com
	BaseUrl string `json:"base_url"`
	// The canonical url of the profile page. For example https://example.com
	ProfileCanonicalUrl string `json:"profile_canonical_url"`
	// The password to authenticate
	Password string `json:"password"`
	// A random string to sign the jwt tokens. Should be at least 32 characters long
	JWTSecret string `json:"jwt_secret"`
	// The store module to use
	StoreData config.ModuleRaw `json:"store" config:"indieauth.store"`
}

func init() {
	config.RegisterModule(&indieAuthPlugin{})
}

func (p *indieAuthPlugin) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "indieauth",
		New:  func() config.Module { return new(indieAuthPlugin) },
	}
}

func (p *indieAuthPlugin) Load(config config.GlobalConfig, _ interface{}) (config.ModuleInstance, error) {

	// remove trailing slash of baseUrl
	p.BaseUrl = strings.TrimSuffix(p.BaseUrl, "/")

	storeInt, err := config.Config.LoadModule(p, "StoreData", nil)
	if err != nil {
		return nil, err
	}
	store, ok := storeInt.(Store)
	if !ok {
		return nil, fmt.Errorf("store module is not of type indieauth.Store: %T", storeInt)
	}

	return NewIndieAuthApiModule(
		p.BaseUrl,
		p.ProfileCanonicalUrl,
		p.Password,
		p.JWTSecret,
		store,
		*config.HttpClient,
	), nil
}
