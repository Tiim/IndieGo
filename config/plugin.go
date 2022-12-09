package config

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type ModuleRaw struct {
	Name string          `json:"name"`
	Args json.RawMessage `json:"args"`
}

type ModuleInfo struct {
	Name string
	New  func() Module
}

type Module interface {
	Load(config GlobalConfig, args interface{}) (ModuleInstance, error)
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
