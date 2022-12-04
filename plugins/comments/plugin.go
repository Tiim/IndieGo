package comments

import (
	"encoding/json"
	"fmt"
	"log"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/plugins/admin"
)

type commentsPlugin struct{}

type commentsPluginData struct {
	StoreData config.ModuleRaw `json:"store"`
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
	storeInt, err := config.Config.LoadModule(d.StoreData)
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

	return NewCommentModule(store), nil
}
