package main

import (
	"log"
	"tiim/go-comment-api/api"
	"tiim/go-comment-api/plugin"

	_ "tiim/go-comment-api/indieauth"
)

func main() {

	configPath := "config.json"
	configStr, err := plugin.ReadConfigString(configPath)
	if err != nil {
		log.Fatalf("unable to read config file: %v", err)
	}
	config, err := plugin.LoadConfig(configStr)
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	apiServer := api.NewApiServer(config.Plugins)
	r, err := apiServer.Start()
	if err != nil {
		log.Fatalf("unable to start api server: %v", err)
	}
	config.StartPlugins()
	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
}
