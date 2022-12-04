package config

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type ModuleRaw struct {
	Name string          `json:"name"`
	Args json.RawMessage `json:"args"`
}

type Plugin interface {
	Load(data json.RawMessage, config GlobalConfig) (PluginInstance, error)
	Name() string
}

type Module interface {
	Load(data json.RawMessage, config GlobalConfig, args interface{}) (ModuleInstance, error)
	Name() string
}

type PluginInstance interface {
	Name() string
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

type ModuleInstance interface{}
