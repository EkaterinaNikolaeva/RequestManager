package main

import (
	"log"
	"os"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/config"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/mattermostprovider"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/messagesmatcher"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/service"
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
	provider := mattermostprovider.NewMattermostProvider(mattermostBot)
	matcher := messagesmatcher.MessagesMatcher{}
	go provider.Run()
	taskFromMessagesCreator := service.NewTaskFromMessagesCreator(provider, matcher)
	taskFromMessagesCreator.Run()
}
