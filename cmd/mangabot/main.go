package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	vkbot "github.com/Kvertinum01/mangabot/internal/app/bot"
)

var (
	configPath string
)

func init() {
	// Parse path to config from command args
	flag.StringVar(&configPath, "config-path", "configs/mangabot.toml", "path to config file")
}

func main() {
	flag.Parse()

	// Read .toml and write to config
	config := &vkbot.Conifg{}
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	// Setup bot
	if err := vkbot.SetupBot(config); err != nil {
		log.Fatal(err)
	}
}
