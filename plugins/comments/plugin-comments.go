package comments

import (
	"encoding/json"
	"fmt"
	"log"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/plugins/admin"
	"tiim/go-comment-api/plugins/shared-modules/event"
)

type commentsPlugin struct{}

type commentsPluginData struct {
	StoreData    config.ModuleRaw `json:"store"`
	EventHandler config.ModuleRaw `json:"event_handler"`
}

func init() {
	config.RegisterPlugin(&commentsPlugin{})
}

func (p *commentsPlugin) Name() string {
	return "comments"
}

func (p *commentsPlugin) Load(data json.RawMessage, config config.GlobalConfig) (config.PluginInstance, error) {
	var d commentsPluginData
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}
	storeInt, err := config.Config.LoadModule(d.StoreData, nil)
	if err != nil {
		return nil, err
	}
	store, ok := storeInt.(commentStore)
	if !ok {
		return nil, fmt.Errorf("store module is not of type comments.commentStore: %T", storeInt)
	}

	adminInt, err := config.GetPlugin("admin")
	if err == nil {
		admin, ok := adminInt.(*admin.AdminModule)
		if !ok {
			return nil, fmt.Errorf("admin is not a of type admin.AdminModule: %T", adminInt)
		}
		admin.RegisterSection(newAdminCommentSection(store))
	} else {
		log.Printf("comments plugin: admin plugin not loaded, not registering admin section")
	}

	eventHandlerInt, err := config.Config.LoadModule(d.EventHandler, store)
	if err != nil {
		return nil, fmt.Errorf("error loading event handler: %v", err)
	}
	eventHandler, ok := eventHandlerInt.(event.Handler)
	if !ok {
		return nil, fmt.Errorf("comments-event-handler is not a of type event.Handler: %T", eventHandlerInt)
	}
	store.SetEventHandler(eventHandler)

	return NewCommentModule(store), nil
}
