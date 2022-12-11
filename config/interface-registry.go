package config

import "log"

var interfaces = make(map[string][]interface{})

func RegisterInterface(name string) {
	if _, ok := interfaces[name]; !ok {
		log.Printf("Registering interface '%s'\n", name)
		interfaces[name] = make([]interface{}, 0)
	}
}

func addInterface(name string, iface interface{}) {
	interfaces[name] = append(interfaces[name], iface)
}

func getInterfaces(name string) []interface{} {
	return interfaces[name]
}
