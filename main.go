package main

import (
	"log"
	"tiim/go-comment-api/api"
	"tiim/go-comment-api/config"

	_ "tiim/go-comment-api/plugins/admin"
	_ "tiim/go-comment-api/plugins/comments"
	_ "tiim/go-comment-api/plugins/indieauth"
	_ "tiim/go-comment-api/plugins/micropub"
	_ "tiim/go-comment-api/plugins/wmreceive"
	_ "tiim/go-comment-api/plugins/wmsend"
)

func main() {

	configPath := "config.json"
	configStr, err := config.ReadConfigString(configPath)
	if err != nil {
		log.Fatalf("unable to read config file: %v", err)
	}
	config, err := config.LoadConfig(configStr)
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	err = config.Init()
	if err != nil {
		log.Fatalf("unable to init config: %v", err)
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
