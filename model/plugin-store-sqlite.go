package model

import (
	"log"
	"tiim/go-comment-api/config"
	"time"
)

type sqliteStoreModule struct{}

func init() {
	config.RegisterModule(&sqliteStoreModule{})
}

func (p *sqliteStoreModule) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "store.sqlite",
		New:  func() config.Module { return new(sqliteStoreModule) },
	}
}

func (p *sqliteStoreModule) Load(config config.GlobalConfig, _ interface{}) (config.ModuleInstance, error) {
	return NewSQLiteStore(config.Scheduler)
}

func (p *SQLiteStore) Name() string {
	return "store.sqlite"
}
func (p *SQLiteStore) Init(config config.GlobalConfig) error {
	p.runMigrations()
	return nil
}
func (p *SQLiteStore) Start() error {
	// TODO: expose cleanup interval as config option
	p.scheduler.Every(12 * time.Hour).Do(func() {
		err := p.CleanUp()
		if err != nil {
			log.Printf("unable to clean up sqlite database: %v", err)
		}
	})
	return nil
}
