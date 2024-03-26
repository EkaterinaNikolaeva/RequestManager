package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/jirahttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/mattermosthttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/config"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/jirataskcreator"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/mattermostprovider"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/mattermostsender"
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	mattermostBot := bot.NewMattermostBot(config)
	jiraHttpClient := jirahttpclient.NewJiraHttpClient(&http.Client{}, config.JiraBaseUrl, config.JiraBotUsername, config.JiraBotPassword)
	jiraTaskCreator := jirataskcreator.NewJiraTaskCreator(jiraHttpClient)
	httpClientForMessanger := mattermosthttpclient.NewHttpClient(&http.Client{}, mattermostBot.Token, config.MattermostHttp)
	provider := mattermostprovider.NewMattermostProvider(mattermostBot)
	sender := mattermostsender.NewMattermostSender(httpClientForMessanger)
	matcher, err := messagesmatcher.NewMessagesMatcher(config.MessagesPattern)
	if err != nil {
		log.Fatalf("Unsuccessful start: %q", err)
	}
	go provider.Run(ctx)
	taskFromMessagesCreator := service.NewTaskFromMessagesCreator(provider, sender, matcher, jiraTaskCreator,
		config.MessageReply, config.JiraProject, config.JiraIssueType)
	taskFromMessagesCreator.Run(ctx)
}
