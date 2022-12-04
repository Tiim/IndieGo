package config

import (
	"fmt"
	"log"
)

// plugins are top level components that provide some functionality,
// e.g. api endpoints, background tasks, ui pages etc.
var plugins = make(map[string]Plugin)

// modules are components that are used by plugins to
// customize their behavior, e.g. a database connection,
// a cache, notification service etc.
var modules = make(map[string]Module)

func RegisterPlugin(a Plugin) {
	log.Printf("Registering plugin %s", a.Name())
	if _, ok := plugins[a.Name()]; ok {
		panic(fmt.Sprintf("plugin %s already registered", a.Name()))
	}
	plugins[a.Name()] = a
}

func RegisterModule(m Module) {
	log.Printf("Registering module %s", m.Name())
	if _, ok := plugins[m.Name()]; ok {
		panic(fmt.Sprintf("module %s already registered", m.Name()))
	}
	modules[m.Name()] = m
}

func (c *Config) LoadPlugins() error {
	for i, pluginRaw := range c.PluginsRaw {
		p, ok := plugins[pluginRaw.Name]
		log.Printf("Loading plugin: '%s'\n", pluginRaw.Name)
		if !ok {
			return fmt.Errorf("plugin '%s' not found", pluginRaw.Name)
		}
		plugin, err := p.Load(pluginRaw.Args, c.GlobalConfig)
		if err != nil {
			return fmt.Errorf("failed to load plugin '%s' (%d): %w", pluginRaw.Name, i, err)
		}
		log.Printf("Loaded  plugin: '%s'\n", plugin.Name())
		if plugin.Name() != pluginRaw.Name {
			return fmt.Errorf("plugin '%s' (%d) returned wrong name: %s", pluginRaw.Name, i, plugin.Name())
		}
		c.Plugins = append(c.Plugins, plugin)
	}
	log.Printf("Successfully loaded %d plugins and %d modules\n", len(c.Plugins), len(c.Modules))
	return nil
}

func (c *Config) LoadModule(moduleData ModuleRaw, args interface{}) (ModuleInstance, error) {
	log.Printf("Loading module: '%s'\n", moduleData.Name)
	m, ok := modules[moduleData.Name]
	if !ok {
		return nil, fmt.Errorf("module %s not found", moduleData.Name)
	}
	module, err := m.Load(moduleData.Args, c.GlobalConfig, args)
	c.Modules = append(c.Modules, module)
	return module, err
}
