package indieauth

import (
	"fmt"
	"log"
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
		Docs: config.ConfigDocs{
			DocString: `IndieAuth module. This module enables the IndieAuth authentication.`,
			Fields: map[string]string{
				"BaseUrl":             "The url indiego is running on. For example https://indiego.example.com",
				"ProfileCanonicalUrl": "The canonical url of the profile page. For example https://example.com",
				"Password":            "The password to authenticate",
				"JWTSecret":           "A random string to sign the jwt tokens. Should be at least 32 characters long",
				"StoreData":           "The store module to use",
			},
		},
	}
}

func (p *indieAuthPlugin) Load(config config.GlobalConfig, _ interface{}, logger *log.Logger) (config.ModuleInstance, error) {

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
		logger,
	), nil
}
