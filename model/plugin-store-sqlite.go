package model

import (
	"encoding/json"
	"log"
	"tiim/go-comment-api/config"
	"time"
)

type sqliteStorePlugin struct{}

func init() {
	config.RegisterPlugin(&sqliteStorePlugin{})
}

func (p *sqliteStorePlugin) Name() string {
	return "store-sqlite"
}

func (p *sqliteStorePlugin) Load(data json.RawMessage, config config.GlobalConfig) (config.PluginInstance, error) {
	return NewSQLiteStore(config.Scheduler)
}

func (p *SQLiteStore) Name() string {
	return "store-sqlite"
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
