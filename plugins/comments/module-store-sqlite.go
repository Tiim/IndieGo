package comments

import (
	"fmt"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/model"
)

type commentSQLiteStoreModule struct {
	// The page mapper module to use for mapping comment ids to urls.
	PageMapper config.ModuleRaw `json:"page_mapper" config:"comments.page-mapper"`
}

func init() {
	config.RegisterModule(&commentSQLiteStoreModule{})
}

func (m *commentSQLiteStoreModule) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "comments.store.sqlite",
		New:  func() config.Module { return new(commentSQLiteStoreModule) },
	}
}

func (m *commentSQLiteStoreModule) Load(config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	storeInt, err := config.GetModule("store.sqlite")
	if err != nil {
		return nil, fmt.Errorf("depends on store.sqlite plugin: %v", err)
	}
	store, ok := storeInt.(*model.SQLiteStore)
	if !ok {
		return nil, fmt.Errorf("store.sqlite is not a of type model.SQLiteStore: %T", storeInt)
	}
	pageMapperInt, err := config.Config.LoadModule(m, "PageMapper", nil)
	if err != nil {
		return nil, fmt.Errorf("error loading page mapper: %v", err)
	}
	pageMapper, ok := pageMapperInt.(CommentPageToUrlMapper)
	if !ok {
		return nil, fmt.Errorf("comments-page-mapper is not a of type comments.CommentPageToUrlMapper: %T", pageMapperInt)
	}
	sqliteStore := &commentSQLiteStore{db: store.GetDBConnection(), pageToUrlMapper: pageMapper}
	return sqliteStore, nil
}
