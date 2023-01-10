package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"tiim/go-comment-api/api"
	"tiim/go-comment-api/config"

	_ "tiim/go-comment-api/plugins/admin"
	_ "tiim/go-comment-api/plugins/comment-provider"
	_ "tiim/go-comment-api/plugins/comments"
	_ "tiim/go-comment-api/plugins/indieauth"
	_ "tiim/go-comment-api/plugins/manual-backup"
	_ "tiim/go-comment-api/plugins/micropub"
	_ "tiim/go-comment-api/plugins/public-site"
	_ "tiim/go-comment-api/plugins/wmreceive"
	_ "tiim/go-comment-api/plugins/wmsend"
)

func main() {

	ensureTempDir()
	var configPath string

	flag.StringVar(&configPath, "config", "config.json", "path to config file")
	genDocs := flag.Bool("generate-docs", false, "generate documentation")
	flag.Parse()

	if *genDocs {
		generateDocs()
	}

	configStr, err := config.ReadConfigString(configPath)
	if err != nil {
		log.Fatalf("unable to read config file: %v", err)
	}

	configStr = os.ExpandEnv(configStr)

	config, err := config.LoadConfig(configStr)
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	err = config.Init()
	if err != nil {
		log.Fatalf("unable to init config: %v", err)
	}

	apiServer := api.NewApiServer(config.Modules)
	r, err := apiServer.Start()
	if err != nil {
		log.Fatalf("unable to start api server: %v", err)
	}
	config.StartModules()
	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
}

func generateDocs() {
	docs := config.GenerateDocs()
	file := "docs.html"
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("unable to create file %s: %v", file, err)
	}
	defer f.Close()
	_, err = f.WriteString(docs)
	if err != nil {
		log.Fatalf("unable to write to file %s: %v", file, err)
	}
	os.Exit(0)
}

// make sure we have a working tempdir, because:
// os.TempDir(): The directory is neither guaranteed to exist nor have accessible permissions.
// https://blog.cubieserver.de/2020/go-debugging-why-parsemultipartform-returns-error-no-such-file-or-directory/
func ensureTempDir() {
	logger := log.New(os.Stdout, "[init] ", log.Flags())
	tempDir := os.TempDir()
	if err := os.MkdirAll(tempDir, 1777); err != nil {
		logger.Fatalf("Failed to create temporary directory %s: %s", tempDir, err)
	}
	tempFile, err := ioutil.TempFile("", "genericInit_")
	if err != nil {
		logger.Fatalf("Failed to create tempFile: %s", err)
	}
	_, err = fmt.Fprintf(tempFile, "Hello, World!")
	if err != nil {
		logger.Fatalf("Failed to write to tempFile: %s", err)
	}
	if err := tempFile.Close(); err != nil {
		logger.Fatalf("Failed to close tempFile: %s", err)
	}
	if err := os.Remove(tempFile.Name()); err != nil {
		logger.Fatalf("Failed to delete tempFile: %s", err)
	}
	logger.Printf("Using temporary directory %s", tempDir)
}
