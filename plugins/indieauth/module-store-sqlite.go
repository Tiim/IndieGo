package indieauth

import (
	"fmt"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/model"
	"time"
)

type indieAuthSQLiteStoreModule struct {
	// The expiration time of auth codes in minutes.
	// The client must register an auth token within this time.
	// Default: 10
	AuthCodeExpirationMinutes int `json:"auth_code_expiration_min"`
	// The expiration time of auth tokens in minutes.
	// The client must re authenticate after this time.
	// Default: 60 * 24 * 30 (30 days)
	AuthTokenExpirationMinutes int `json:"auth_token_expiration_min"`
}

func init() {
	config.RegisterModule(&indieAuthSQLiteStoreModule{})
}

func (m *indieAuthSQLiteStoreModule) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "indieauth.store.sqlite",
		New: func() config.Module {
			m := new(indieAuthSQLiteStoreModule)
			m.AuthCodeExpirationMinutes = 10
			m.AuthTokenExpirationMinutes = 60 * 24 * 30
			return m
		},
		Docs: config.ConfigDocs{
			DocString: `SQLite store module. This module stores the auth codes and tokens in a SQLite database.`,
			Fields: map[string]string{
				"AuthCodeExpirationMinutes":  "The expiration time of auth codes in minutes. The client must register an auth token within this time. Default: 10",
				"AuthTokenExpirationMinutes": "The expiration time of auth tokens in minutes. The client must re authenticate after this time. Default: 60 * 24 * 30 (30 days)",
			},
		},
	}
}

func (m *indieAuthSQLiteStoreModule) Load(config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	storeInt, err := config.GetModule("store.sqlite")
	if err != nil {
		return nil, fmt.Errorf("depends on store.sqlite plugin: %v", err)
	}
	store, ok := storeInt.(*model.SQLiteStore)
	if !ok {
		return nil, fmt.Errorf("store.sqlite is not a of type model.SQLiteStore: %T", storeInt)
	}
	return NewSQLiteStore(
		store.GetDBConnection(),
		time.Duration(m.AuthCodeExpirationMinutes)*time.Minute,
		time.Duration(m.AuthTokenExpirationMinutes)*time.Minute,
	), nil
}
