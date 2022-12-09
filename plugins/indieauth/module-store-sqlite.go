package indieauth

import (
	"encoding/json"
	"fmt"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/model"
	"time"
)

type indieAuthSQLiteStoreModule struct{}
type indieAuthSQLiteStoreModuleData struct {
	AuthCodeExpirationMinutes  int `json:"authCodeExpirationMinutes"`
	AuthTokenExpirationMinutes int `json:"authTokenExpirationMinutes"`
}

func init() {
	config.RegisterModule(&indieAuthSQLiteStoreModule{})
}

func (m *indieAuthSQLiteStoreModule) Name() string {
	return "indieauth-store-sqlite"
}

func (m *indieAuthSQLiteStoreModule) Load(data json.RawMessage, config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	d := indieAuthSQLiteStoreModuleData{
		AuthCodeExpirationMinutes:  10,
		AuthTokenExpirationMinutes: 60 * 24 * 30,
	}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}
	storeInt, err := config.GetPlugin("store-sqlite")
	if err != nil {
		return nil, fmt.Errorf("%s depends on store-sqlite plugin, error loading: %v", m.Name(), err)
	}
	store, ok := storeInt.(*model.SQLiteStore)
	if !ok {
		return nil, fmt.Errorf("store-sqlite is not a of type model.SQLiteStore: %T", storeInt)
	}
	return NewSQLiteStore(
		store.GetDBConnection(),
		time.Duration(d.AuthCodeExpirationMinutes)*time.Minute,
		time.Duration(d.AuthTokenExpirationMinutes)*time.Minute,
	), nil
}
