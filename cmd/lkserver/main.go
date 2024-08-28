package main

import (
	"flag"
	"lkserver/internal/lkserver"
	"lkserver/internal/repository"
	"lkserver/internal/repository/json"
	"log"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "res/config.toml", "Path co config file")
	flag.Parse()
}

func main() {
	config := lkserver.NewConfig()
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		log.Fatal(err)
	}

	repo, err := initRepo()
	if err != nil {
		log.Fatal("Error on init repository")
	}
	defer repo.Close()

	server := lkserver.New(
		repo,
		config,
	)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Bind addr: %s", config.BindAddr)
}

func initRepo() (*repository.Repo, error) {

	userRepo, err := json.NewUserRepo("data")
	if err != nil {
		return nil, err
	}

	return &repository.Repo{
		User: userRepo,
	}, nil
}
