package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	ModulesRaw []ModuleRaw               `json:"modules"`
	Modules    map[string]ModuleInstance `json:"-"`
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

func (c *Config) Init() error {
	log.Println("Initializing modules")
	for name, module := range c.Modules {
		if plugin, ok := module.(PluginInstance); ok {
			err := plugin.Init(c.GlobalConfig)
			if err != nil {
				return fmt.Errorf("failed initializing plugin %s: %v", name, err)
			}
		}
	}
	log.Println("Initializing modules done")
	return nil
}

func (c *Config) StartModules() error {
	for name, module := range c.Modules {
		if plugin, ok := module.(PluginInstance); ok {
			err := plugin.Start()
			if err != nil {
				return fmt.Errorf("failed starting plugin %s: %s", name, err)
			}
		}
	}
	return nil
}

func (gc *GlobalConfig) GetModule(name string) (ModuleInstance, error) {
	module, ok := gc.Config.Modules[name]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", name)
	}
	return module, nil
}
