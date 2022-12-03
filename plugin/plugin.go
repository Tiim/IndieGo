package plugin

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
	Load(data json.RawMessage, config GlobalConfig) (ModuleInstance, error)
	Name() string
}

type PluginInstance interface {
	Name() string
	Init(r *gin.Engine) error
	RegisterRoutes(r *gin.Engine) error
	Start() error
}

type ModuleInstance interface{}
