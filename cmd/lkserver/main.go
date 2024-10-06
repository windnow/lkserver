package main

import (
	"flag"
	"lkserver/internal/lkserver"
	"lkserver/internal/lkserver/config"
	"lkserver/internal/repository/file"
	"lkserver/internal/repository/sqlite"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "res/config.toml", "Path co config file")
	flag.Parse()
}

func main() {
	config, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// repo, err := json.NewJSONProvider("data/json")
	repo, err := sqlite.NewSQLiteProvider("data/data.db")
	if err != nil {
		log.Fatal("Error on init repository: ", err.Error())
	}
	defer repo.Close()
	fileRepo, err := file.NewFileRepo("data/files")
	if err != nil {
		log.Fatal("Error on init file repository: ", err.Error())
	}

	server := lkserver.New(
		repo,
		fileRepo,
		config,
	)

	if err := server.Start(); err != nil {
		log.Fatal("Error starting repository: ", err)
	}

	log.Printf("Bind addr: %s", config.BindAddr)
}
