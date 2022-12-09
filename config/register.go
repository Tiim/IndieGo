package config

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
)

// modules are components that are used by plugins to
// customize their behavior, e.g. a database connection,
// a cache, notification service etc.
var modules = make(map[string]Module)

func RegisterModule(a Module) {
	info := a.IndieGoModule()
	log.Printf("Registering plugin %s", info.Name)
	if _, ok := modules[info.Name]; ok {
		panic(fmt.Sprintf("plugin %s already registered", info.Name))
	}
	modules[info.Name] = a
}

func (c *Config) LoadPlugins() error {
	c.Modules = make(map[string]ModuleInstance)
	for i, moduleRaw := range c.ModulesRaw {
		name := moduleRaw.Name
		p, ok := modules[name]
		log.Printf("Loading plugin: '%s'\n", name)
		if !ok {
			return fmt.Errorf("plugin '%s' not found", name)
		}

		moduleInstance := p.IndieGoModule().New()

		if moduleRaw.Args != nil {
			err := json.Unmarshal(moduleRaw.Args, moduleInstance)
			if err != nil {
				return fmt.Errorf("failed to unmarshal plugin '%s' (%d): %w", name, i, err)
			}
		}

		module, err := moduleInstance.Load(c.GlobalConfig, nil)
		if err != nil {
			return fmt.Errorf("failed to load plugin '%s' (%d): %w", name, i, err)
		}
		c.Modules[name] = module
	}
	log.Printf("Successfully loaded %d modules\n", len(c.Modules))
	return nil
}

func (c *Config) LoadModule(structPtr any, fieldName string, args any) (ModuleInstance, error) {
	mi, err := c.loadModule(structPtr, fieldName, args)
	if err != nil {
		return nil, err
	}
	module, ok := mi.(ModuleInstance)
	if !ok {
		return nil, fmt.Errorf("field %s is not a single module", fieldName)
	}
	return module, nil
}

func (c *Config) LoadModuleSlice(structPtr any, fieldName string, args any) ([]ModuleInstance, error) {
	mi, err := c.loadModule(structPtr, fieldName, args)
	if err != nil {
		return nil, err
	}
	modules, ok := mi.([]ModuleInstance)
	if !ok {
		return nil, fmt.Errorf("field %s is not a slice of modules", fieldName)
	}
	return modules, nil
}

func (c *Config) loadModule(structPtr any, fieldName string, args interface{}) (any, error) {

	val := reflect.ValueOf(structPtr).Elem().FieldByName(fieldName)

	field, ok := reflect.TypeOf(structPtr).Elem().FieldByName(fieldName)
	if !ok {
		return nil, fmt.Errorf("field %s does not exist in %T", fieldName, structPtr)
	}

	nameSpace := field.Tag.Get("config")

	if c.Modules == nil {
		c.Modules = make(map[string]ModuleInstance)
	}
	switch t := val.Interface().(type) {
	case ModuleRaw:

		if t.Name == "" {
			return nil, fmt.Errorf("module has no name: %s (%s) child of %#v", fieldName, nameSpace, structPtr)
		}

		module, err := c.loadSingleModule(t, nameSpace, fieldName)
		if err != nil {
			return nil, fmt.Errorf("failed to load module %s (field %s): %w", fieldName, fieldName, err)
		}
		moduleInstance, err := module.Load(c.GlobalConfig, args)
		if err != nil {
			return nil, fmt.Errorf("failed to load module %s (field %s): %w", fieldName, fieldName, err)
		}
		c.Modules[t.Name] = moduleInstance
		return moduleInstance, nil
	case []ModuleRaw:
		modules := make([]ModuleInstance, len(t))
		for i, moduleData := range t {
			if moduleData.Name == "" {
				return nil, fmt.Errorf("module has no name: %s (%s) child of %#v", fieldName, nameSpace, structPtr)
			}
			module, err := c.loadSingleModule(moduleData, nameSpace, fieldName)
			if err != nil {
				return nil, fmt.Errorf("failed to load module %s (field %s, index %d): %w", fieldName, fieldName, i, err)
			}
			moduleInstance, err := module.Load(c.GlobalConfig, args)
			if err != nil {
				return nil, fmt.Errorf("failed to load module %s (field %s, index %d): %w", fieldName, fieldName, i, err)
			}
			c.Modules[moduleData.Name] = moduleInstance
			modules[i] = moduleInstance
		}
		return modules, nil
	default:
		return nil, fmt.Errorf("unknown type for field %s: %T, must be config.ModuleRaw or []config.ModuleRaw", fieldName, val.Interface())
	}

}

func (c *Config) loadSingleModule(moduleData ModuleRaw, nameSpace, fieldName string) (Module, error) {
	name := moduleData.Name

	actualNamespace := GetNamespaceFromName(name)

	if actualNamespace != nameSpace {
		return nil, fmt.Errorf("field %s has namespace '%s' but module %s has namespace '%s'", fieldName, nameSpace, name, actualNamespace)
	}

	log.Printf("Loading module: '%s'\n", name)
	m, ok := modules[name]
	if !ok {
		return nil, fmt.Errorf("module %s not found", name)
	}

	module := m.IndieGoModule().New()
	if moduleData.Args != nil {
		err := json.Unmarshal(moduleData.Args, module)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal module %s: %w", name, err)
		}
	}
	return module, nil
}

func GetNamespaceFromName(name string) string {
	namespaceSegments := strings.Split(name, ".")
	return strings.Join(namespaceSegments[:len(namespaceSegments)-1], ".")
}
