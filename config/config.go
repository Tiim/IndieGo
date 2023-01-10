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
	PluginsRaw []ModuleRaw `json:"plugins"`

	Modules map[string][]ModuleInstance `json:"-"`
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
		return nil, fmt.Errorf("unable to unmarshal top level json: %w", err)
	}
	config.GlobalConfig.Config = config

	config.GlobalConfig.HttpClient = &http.Client{Timeout: time.Second * 10}
	config.GlobalConfig.Scheduler = gocron.NewScheduler(time.UTC)

	err = config.LoadPlugins()

	if err != nil {
		return nil, fmt.Errorf("unable to load plugins: %w", err)
	}

	return config, nil
}

func (c *Config) Init() error {
	log.Println("Initializing modules")
	for name, modules := range c.Modules {
		for _, module := range modules {
			if plugin, ok := module.(PluginInstance); ok {
				err := plugin.Init(c.GlobalConfig)
				if err != nil {
					return fmt.Errorf("failed initializing plugin %s: %v", name, err)
				}
			}
		}
	}
	log.Println("Initializing modules done")
	return nil
}

func (c *Config) StartModules() error {
	for name, modules := range c.Modules {
		for _, module := range modules {
			if plugin, ok := module.(PluginInstance); ok {
				err := plugin.Start()
				if err != nil {
					return fmt.Errorf("failed starting plugin %s: %s", name, err)
				}
			}
		}
	}
	return nil
}

func (gc *GlobalConfig) GetModule(name string) (ModuleInstance, error) {
	module, ok := gc.Config.Modules[name]
	if ok {
		if len(module) > 0 {
			return module[0], nil
		}
	}
	for mname, module := range gc.Config.Modules {
		if GetNamespaceFromName(mname) == name && len(module) > 0 {
			return module[0], nil
		}
	}
	return nil, fmt.Errorf("plugin '%s' not found", name)
}

func (c *Config) AddInterface(name string, iface any) {
	addInterface(name, iface)
}

func (c *Config) GetInterfaces(name string) []interface{} {
	return getInterfaces(name)
}
