package comments

import (
	"encoding/json"
	"fmt"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/model"
)

type commentSQLiteStoreModule struct{}
type commentSQLiteStoreModuelData struct {
	PageMapper config.ModuleRaw `json:"page_mapper"`
}

func init() {
	config.RegisterModule(&commentSQLiteStoreModule{})
}

func (m *commentSQLiteStoreModule) Name() string {
	return "comments-sqlite-store"
}

func (m *commentSQLiteStoreModule) Load(data json.RawMessage, config config.GlobalConfig) (config.ModuleInstance, error) {
	d := commentSQLiteStoreModuelData{}
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
	pageMapperInt, err := config.Config.LoadModule(d.PageMapper)
	if err != nil {
		return nil, fmt.Errorf("error loading page mapper: %v", err)
	}
	pageMapper, ok := pageMapperInt.(CommentPageToUrlMapper)
	if !ok {
		return nil, fmt.Errorf("comments-page-mapper is not a of type comments.CommentPageToUrlMapper: %T", pageMapperInt)
	}
	return &commentSQLiteStore{db: store.GetDBConnection(), pageToUrlMapper: pageMapper}, nil
}
