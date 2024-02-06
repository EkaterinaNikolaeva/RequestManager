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
	config := config.LoadConfig(os.Args[1])
	mattermostBot := bot.NewMattermostBot(config)
	server.MakeServer(mattermostBot)
}
