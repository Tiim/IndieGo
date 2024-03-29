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
		Docs: config.ConfigDocs{
			DocString: `Root SQLite store module. Most sqlite store modules require this module to be loaded first. 
		This module provides the database connection and migrations.`,
			Fields: map[string]string{},
		},
	}
}

func (p *sqliteStoreModule) Load(config config.GlobalConfig, _ interface{}, logger *log.Logger) (config.ModuleInstance, error) {
	return NewSQLiteStore(config.Scheduler, logger)
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
			p.logger.Printf("unable to clean up sqlite database: %v", err)
		}
	})
	return nil
}
