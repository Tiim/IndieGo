package config

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

type ModuleRaw struct {
	Name string          `json:"name"`
	Args json.RawMessage `json:"args"`
}

type ModuleInfo struct {
	Name string
	New  func() Module
	Docs ConfigDocs
}

type Module interface {
	Load(config GlobalConfig, args interface{}, logger *log.Logger) (ModuleInstance, error)
	IndieGoModule() ModuleInfo
}

type ModuleInstance interface{}

type PluginInstance interface {
	ModuleInstance
	Init(config GlobalConfig) error
	Start() error
}

type ApiPluginInstance interface {
	PluginInstance
	RegisterRoutes(r *gin.Engine) error
}

type GroupedApiPluginInstance interface {
	ApiPluginInstance
	InitGroups(r *gin.Engine) error
}
