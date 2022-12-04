package event

import (
	"encoding/json"
	"tiim/go-comment-api/config"
)

type HandlerModule struct{}
type HandlerModuleData struct {
	Handlers []config.ModuleRaw `json:"handlers"`
}

func init() {
	config.RegisterModule(&HandlerModule{})
}

func (m *HandlerModule) Name() string {
	return "event-handler-list"
}

func (m *HandlerModule) Load(data json.RawMessage, config config.GlobalConfig, args interface{}) (config.ModuleInstance, error) {
	d := HandlerModuleData{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}
	handlers := make([]Handler, len(d.Handlers))
	for i, handlerRaw := range d.Handlers {
		handler, err := config.Config.LoadModule(handlerRaw, args)
		if err != nil {
			return nil, err
		}
		handlers[i] = handler.(Handler)
	}
	return &handlerList{handlers: handlers}, nil
}
