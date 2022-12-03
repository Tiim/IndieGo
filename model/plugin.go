package model

import (
	"encoding/json"
	"log"
	"tiim/go-comment-api/plugin"
	"time"

	"github.com/gin-gonic/gin"
)

type sqliteStorePlugin struct{}

func init() {
	plugin.RegisterPlugin(&sqliteStorePlugin{})
}

func (p *sqliteStorePlugin) Name() string {
	return "sqlite-store"
}

func (p *sqliteStorePlugin) Load(data json.RawMessage, config plugin.GlobalConfig) (plugin.PluginInstance, error) {
	return NewSQLiteStore(config.Scheduler)
}

func (p *SQLiteStore) Name() string {
	return "sqlite-store"
}
func (p *SQLiteStore) Init(r *gin.Engine) error {
	p.runMigrations()
	return nil
}
func (p *SQLiteStore) RegisterRoutes(r *gin.Engine) error {
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
