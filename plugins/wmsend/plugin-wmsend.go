package wmsend

import (
	"fmt"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/plugins/shared-modules/trigger"
	"time"
)

type wmSendPlugin struct {
	// FeedUrl is the URL of the RSS feed. This feed gets periodically polled
	// for new entries. When a new entry is found, webmentions get sent to
	// all URLs found in the entry.
	FeedUrl string `json:"feed_url"`
	// SendIntervalMinutes is the interval in minutes at which the RSS feed
	// gets polled for new entries.
	// Default: 60
	IntervalMinutes int              `json:"interval_minutes"`
	StoreData       config.ModuleRaw `json:"store" config:"webmention.send.store"`
	Trigger         config.ModuleRaw `json:"trigger" config:"trigger"`
}

func init() {
	config.RegisterModule(&wmSendPlugin{})
}

func (p *wmSendPlugin) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "webmention.send",
		New:  func() config.Module { return new(wmSendPlugin) },
	}
}

func (p *wmSendPlugin) Load(config config.GlobalConfig, _ interface{}) (config.ModuleInstance, error) {

	if p.IntervalMinutes == 0 {
		p.IntervalMinutes = 60
	}

	storeInt, err := config.Config.LoadModule(p, "StoreData", nil)
	if err != nil {
		return nil, err
	}
	store := storeInt.(WmSendStore)

	wmModule := &wmSend{
		store:     store,
		rss:       p.FeedUrl,
		client:    config.HttpClient,
		scheduler: config.Scheduler,
		interval:  time.Minute * time.Duration(p.IntervalMinutes),
	}

	var trig trigger.Trigger
	if p.Trigger.Name != "" {
		triggerInt, err := config.Config.LoadModule(p, "Trigger", nil)
		if err != nil {
			return nil, fmt.Errorf("error loading trigger: %v", err)
		}
		trig = triggerInt.(trigger.Trigger)
		trig.AddCallback(func() {
			wmModule.SendNow()
		})
	}

	return wmModule, nil
}
