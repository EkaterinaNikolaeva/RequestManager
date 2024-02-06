package main

import (
	"log"
	"os"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/config"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("There is not enough args: file name of config")
	}
	config, err := config.LoadConfig(os.Args[1])
	if err != nil {
		log.Fatalf("Error when opening config file: %q", err)
	}
	mattermostBot := bot.NewMattermostBot(config)
	server.MakeServer(mattermostBot)
}
