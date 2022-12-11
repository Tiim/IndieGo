package comments

import (
	"fmt"
	"log"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/plugins/admin"
	commentprovider "tiim/go-comment-api/plugins/comment-provider"
	"tiim/go-comment-api/plugins/shared-modules/event"
)

type commentsPlugin struct {
	// The store module to use for storing comments.
	// This module must implement the comments.commentStore interface.
	Store config.ModuleRaw `json:"store" config:"comments.store"`
	// The event handler, which will be notified about new comments and
	// deleted comments.
	EventHandler config.ModuleRaw `json:"event_handler" config:"event.mention"`
}

func init() {
	config.RegisterModule(&commentsPlugin{})
}

func (p *commentsPlugin) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "comments",
		New:  func() config.Module { return new(commentsPlugin) },
	}
}

func (p *commentsPlugin) Load(config config.GlobalConfig, _ interface{}) (config.ModuleInstance, error) {
	storeInt, err := config.Config.LoadModule(p, "Store", nil)
	if err != nil {
		return nil, err
	}
	store, ok := storeInt.(commentStore)
	if !ok {
		return nil, fmt.Errorf("store module is not of type comments.commentStore: %T", storeInt)
	}

	adminInt, err := config.GetModule("admin")
	if err == nil {
		admin, ok := adminInt.(*admin.AdminModule)
		if !ok {
			return nil, fmt.Errorf("admin is not a of type admin.AdminModule: %T", adminInt)
		}
		admin.RegisterSection(newAdminCommentSection(store))
	} else {
		log.Printf("comments plugin: admin plugin not loaded, not registering admin section")
	}

	eventHandlerInt, err := config.Config.LoadModule(p, "EventHandler", store)
	if err != nil {
		return nil, fmt.Errorf("error loading event handler: %v", err)
	}
	eventHandler, ok := eventHandlerInt.(event.Handler)
	if !ok {
		return nil, fmt.Errorf("comments-event-handler is not a of type event.Handler: %T", eventHandlerInt)
	}
	store.SetEventHandler(eventHandler)

	var commentProvider commentprovider.CommentProvider = store
	config.Config.AddInterface("comment-provider.provider", commentProvider)

	return NewCommentModule(store), nil
}
