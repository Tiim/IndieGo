package wmsend

import (
	"encoding/json"
	"fmt"
	"tiim/go-comment-api/config"
	"time"
)

type WmSendPlugin struct{}
type WmSendPluginData struct {
	FeedUrl             string           `json:"feed_url"`
	SendIntervalMinutes int              `json:"send_interval_minutes"`
	StoreData           config.ModuleRaw `json:"store"`
}

func init() {
	config.RegisterPlugin(&WmSendPlugin{})
}

func (p *WmSendPlugin) Name() string {
	return "webmention-send"
}

func (p *WmSendPlugin) Load(data json.RawMessage, config config.GlobalConfig) (config.PluginInstance, error) {

	d := WmSendPluginData{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}

	storeInt, err := config.Config.LoadModule(d.StoreData, nil)
	if err != nil {
		return nil, err
	}
	store, ok := storeInt.(WmSendStore)
	if !ok {
		return nil, fmt.Errorf("store module is not of type wmsend.WmSendStore: %T", storeInt)
	}

	return &wmSend{
		store:     store,
		rss:       d.FeedUrl,
		client:    config.HttpClient,
		scheduler: config.Scheduler,
		interval:  time.Minute * time.Duration(d.SendIntervalMinutes),
	}, nil
}
