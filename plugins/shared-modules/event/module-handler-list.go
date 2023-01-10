package event

import (
	"fmt"
	"log"
	"tiim/go-comment-api/config"
)

type HandlerModule struct {
	Handlers []config.ModuleRaw `json:"handlers" config:"event.mention"`
}

func init() {
	config.RegisterModule(&HandlerModule{})
}

func (m *HandlerModule) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "event.mention.handler-list",
		New:  func() config.Module { return new(HandlerModule) },
		Docs: config.ConfigDocs{
			DocString: `Handler list module. 
				This module is a list of handlers that are called when a new comment is submitted. 
				Use this module if you want to run multiple handlers.`,
			Fields: map[string]string{
				"Handlers": "List of handlers.",
			},
		},
	}
}

func (m *HandlerModule) Load(config config.GlobalConfig, args interface{}, logger *log.Logger) (config.ModuleInstance, error) {
	h, err := config.Config.LoadModuleSlice(m, "Handlers", args)
	if err != nil {
		return nil, err
	}
	handlers := make([]Handler, len(h))
	for i, v := range h {
		h, ok := v.(Handler)
		if !ok {
			return nil, fmt.Errorf("event.mention.handler-list: handler %d is not a Handler: %T", i, v)
		}
		handlers[i] = h
	}
	return &handlerList{handlers: handlers, logger: logger}, nil
}
