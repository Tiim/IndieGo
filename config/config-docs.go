package config

import (
	"fmt"
	"reflect"
)

type ConfigDocs struct {
	DocString string
	Fields    map[string]string
}

func validateConfigDocs(plugin Module, configDocs ConfigDocs) error {
	if configDocs.DocString == "" {
		return fmt.Errorf("plugin %s docstring is empty", plugin.IndieGoModule().Name)
	}

	typ := reflect.TypeOf(plugin)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if val, ok := configDocs.Fields[field.Name]; !ok || val == "" {
			return fmt.Errorf("plugin %s docstring missing for %q (%s)", plugin.IndieGoModule().Name, field.Name, jsonTag)
		}
	}
	for name, docstring := range configDocs.Fields {
		field, ok := typ.FieldByName(name)
		if !ok {
			return fmt.Errorf("plugin %s docstring no such field %q", plugin.IndieGoModule().Name, name)
		}
		if docstring == "" {
			return fmt.Errorf("plugin %s docstring for %q is empty", plugin.IndieGoModule().Name, name)
		}
		if field.Tag.Get("json") == "-" {
			return fmt.Errorf("plugin %s docstring field %q is not exported", plugin.IndieGoModule().Name, name)
		}
	}
	return nil
}
