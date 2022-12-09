package wmsend

import (
	"encoding/json"
	"fmt"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/model"
)

type wmSendSQLiteStoreModule struct{}

func init() {
	config.RegisterModule(&wmSendSQLiteStoreModule{})
}

func (m *wmSendSQLiteStoreModule) Name() string {
	return "webmention-send-store-sqlite"
}

func (m *wmSendSQLiteStoreModule) Load(data json.RawMessage, config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	storeInt, err := config.GetPlugin("store-sqlite")
	if err != nil {
		return nil, err
	}
	store, ok := storeInt.(*model.SQLiteStore)
	if !ok {
		return nil, fmt.Errorf("store-sqlite is not a of type model.SQLiteStore: %T", storeInt)
	}
	return newWmSendStore(store.GetDBConnection()), nil
}
