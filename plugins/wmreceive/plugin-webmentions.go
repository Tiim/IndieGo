package wmrecv

import (
	"fmt"
	"log"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/model"
	"tiim/go-comment-api/plugins/admin"
	commentprovider "tiim/go-comment-api/plugins/comment-provider"
	"tiim/go-comment-api/plugins/shared-modules/event"
)

type wmReceivePlugin struct {
	// TargetDomains is a list of domains that are allowed to be the target of a webmention.
	TargetDomains []string         `json:"target_domains"`
	EventHandler  config.ModuleRaw `json:"event_handler" config:"event.mention"`
}

func init() {
	config.RegisterModule(&wmReceivePlugin{})
}

func (p *wmReceivePlugin) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "webmention.receive",
		New:  func() config.Module { return new(wmReceivePlugin) },
	}
}

func (p *wmReceivePlugin) Load(config config.GlobalConfig, _ interface{}) (config.ModuleInstance, error) {

	if len(p.TargetDomains) == 0 {
		return nil, fmt.Errorf("at least one target domain must be specified")
	}

	storeInt, err := config.GetModule("store.sqlite")
	if err != nil {
		return nil, fmt.Errorf("webmention-receive plugin requires store.sqlite plugin: %v", err)
	}
	store, ok := storeInt.(*model.SQLiteStore)
	if !ok {
		return nil, fmt.Errorf("store.sqlite is not a of type model.SQLiteStore: %T", storeInt)
	}
	wmStore := newStore(store)
	wmChecker := newWebmentionChecker([]Checker{
		newTargetChecker(p.TargetDomains...),
		newDomainChecker(wmStore),
		newLinkToTargetChecker(),
		newMicroformatEnricherChecker(),
	})
	wmWorker := newMentionsQueueWorker(wmStore, wmChecker)

	eventHandlerInt, err := config.Config.LoadModule(p, "EventHandler", nil)
	if err != nil {
		return nil, fmt.Errorf("error loading event handler: %v", err)
	}
	eventHandler, ok := eventHandlerInt.(event.Handler)
	if !ok {
		return nil, fmt.Errorf("comments-event-handler is not a of type event.Handler: %T", eventHandlerInt)
	}
	wmStore.SetEventHandler(eventHandler)

	adminInt, err := config.GetModule("admin")
	if err == nil {
		admin, ok := adminInt.(*admin.AdminModule)
		if !ok {
			return nil, fmt.Errorf("admin is not a of type admin.AdminModule: %T", admin)
		}
		admin.RegisterSection(newAdminWebmentionsSection(wmStore))
	} else {
		log.Printf("webmention.receive plugin: admin plugin not loaded, not registering admin section")
	}

	var commentProvider commentprovider.CommentProvider = wmStore
	config.Config.AddInterface("comment-provider.provider", commentProvider)

	return newApi(wmStore, wmWorker, config.Scheduler), nil
}
