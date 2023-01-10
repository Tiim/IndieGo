package wmsend

import (
	"fmt"
	"log"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/model"
)

type wmSendSQLiteStoreModule struct{}

func init() {
	config.RegisterModule(&wmSendSQLiteStoreModule{})
}

func (m *wmSendSQLiteStoreModule) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "webmention.send.store.sqlite",
		New:  func() config.Module { return new(wmSendSQLiteStoreModule) },
		Docs: config.ConfigDocs{
			DocString: `SQLite store module for webmention send. Must be loaded after the store.sqlite module.`,
		},
	}
}

func (m *wmSendSQLiteStoreModule) Load(config config.GlobalConfig, args interface{}, logger *log.Logger) (config.ModuleInstance, error) {
	storeInt, err := config.GetModule("store.sqlite")
	if err != nil {
		return nil, err
	}
	store, ok := storeInt.(*model.SQLiteStore)
	if !ok {
		return nil, fmt.Errorf("store.sqlite is not a of type model.SQLiteStore: %T", storeInt)
	}
	return newWmSendStore(store.GetDBConnection(), logger), nil
}
