package indieauth

import (
	"encoding/json"
	"fmt"
	"tiim/go-comment-api/model"
	"tiim/go-comment-api/plugin"
	"time"
)

type indieAuthSQLiteStoreModule struct{}
type indieAuthSqliteStoreModuleData struct {
	AuthCodeExpirationMinutes  int `json:"authCodeExpirationMinutes"`
	AuthTokenExpirationMinutes int `json:"authTokenExpirationMinutes"`
}

func init() {
	plugin.RegisterModule(&indieAuthSQLiteStoreModule{})
}

func (m *indieAuthSQLiteStoreModule) Name() string {
	return "indieauth-sqlite-store"
}

func (m *indieAuthSQLiteStoreModule) Load(data json.RawMessage, config plugin.GlobalConfig) (plugin.ModuleInstance, error) {
	d := indieAuthSqliteStoreModuleData{
		AuthCodeExpirationMinutes:  10,
		AuthTokenExpirationMinutes: 60 * 24 * 30,
	}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}
	storeInt, err := config.GetPlugin("sqlite-store")
	if err != nil {
		return nil, fmt.Errorf("%s depends on sqlite-store plugin, error loading: %v", m.Name(), err)
	}
	store, ok := storeInt.(*model.SQLiteStore)
	if !ok {
		return nil, fmt.Errorf("sqlite-store is not a of type model.SQLiteStore: %T", storeInt)
	}
	return NewSQLiteStore(
		store.GetDBConnection(),
		time.Duration(d.AuthCodeExpirationMinutes)*time.Minute,
		time.Duration(d.AuthTokenExpirationMinutes)*time.Minute,
	), nil
}
