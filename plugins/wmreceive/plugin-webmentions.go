package wmrecv

import (
	"encoding/json"
	"fmt"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/model"
	"tiim/go-comment-api/plugins/shared-modules/event"
)

type wmReceivePlugin struct{}
type wmReceivePluginData struct {
	TargetDomains []string         `json:"target_domains"`
	EventHandler  config.ModuleRaw `json:"event_handler"`
}

func init() {
	config.RegisterPlugin(&wmReceivePlugin{})
}

func (p *wmReceivePlugin) Name() string {
	return "webmention-receive"
}

func (p *wmReceivePlugin) Load(data json.RawMessage, config config.GlobalConfig) (config.PluginInstance, error) {
	var d wmReceivePluginData
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}

	storeInt, err := config.GetPlugin("store-sqlite")
	if err != nil {
		return nil, fmt.Errorf("webmention-receive plugin requires store-sqlite plugin: %v", err)
	}
	store, ok := storeInt.(*model.SQLiteStore)
	if !ok {
		return nil, fmt.Errorf("store-sqlite is not a of type model.SQLiteStore: %T", storeInt)
	}
	wmStore := newStore(store)
	wmChecker := newWebmentionChecker([]Checker{
		newTargetChecker(d.TargetDomains...),
		newDomainChecker(wmStore),
		newLinkToTargetChecker(),
		newMicroformatEnricherChecker(),
	})
	wmWorker := newMentionsQueueWorker(wmStore, wmChecker)

	eventHandlerInt, err := config.Config.LoadModule(d.EventHandler, nil)
	if err != nil {
		return nil, fmt.Errorf("error loading event handler: %v", err)
	}
	eventHandler, ok := eventHandlerInt.(event.Handler)
	if !ok {
		return nil, fmt.Errorf("comments-event-handler is not a of type event.Handler: %T", eventHandlerInt)
	}
	wmStore.SetEventHandler(eventHandler)

	return newApi(wmStore, wmWorker, config.Scheduler), nil
}
