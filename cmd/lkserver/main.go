package main

import (
	"flag"
	"lkserver/internal/lkserver"
	"lkserver/internal/repo/json"
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

	storage, err := json.New("data")
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	server := lkserver.New(storage, config)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Bind addr: %s", config.BindAddr)
}
