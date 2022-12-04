package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
)

type GlobalConfig struct {
	Config     *Config
	HttpClient *http.Client      `json:"-"`
	Scheduler  *gocron.Scheduler `json:"-"`
	DB         *sql.DB           `json:"-"`
}

type Config struct {
	GlobalConfig
	PluginsRaw []ModuleRaw      `json:"plugins"`
	Plugins    []PluginInstance `json:"-"`
}

func ReadConfigString(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	return string(b), err
}

func LoadConfig(configString string) (*Config, error) {
	config := &Config{}
	err := json.Unmarshal([]byte(configString), config)
	if err != nil {
		return nil, err
	}
	config.GlobalConfig.Config = config

	config.GlobalConfig.HttpClient = &http.Client{Timeout: time.Second * 10}
	config.GlobalConfig.Scheduler = gocron.NewScheduler(time.UTC)

	err = config.LoadPlugins()

	return config, err
}

func (c *Config) StartPlugins() error {
	for _, plugin := range c.Plugins {
		err := plugin.Start()
		if err != nil {
			return fmt.Errorf("failed starting plugin %s: %s", plugin.Name(), err)
		}
	}
	return nil
}

func (gc *GlobalConfig) GetPlugin(name string) (PluginInstance, error) {
	for _, plugin := range gc.Config.Plugins {
		if plugin.Name() == name {
			return plugin, nil
		}
	}
	return nil, fmt.Errorf("plugin '%s' not found", name)
}
