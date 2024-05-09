package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/config"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/messagesmatcher"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/messagesprovider/mattermostprovider"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/messagesprovider/rocketchatprovider"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/messagessender/mattermostsender"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/messagessender/rocketchatsender"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/service"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/storage/storageinmemory"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/storage/storagepostgres"
	"github.com/gopackage/ddp"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/jirahttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/mattermosthttpclient"
	rocketchathttpclient "github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/rocketchatclient"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/jiracommentcreator"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/jirataskcreator"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("There is not enough args: file name of config")
	}
	configData, err := config.LoadConfig(os.Args[1])
	if err != nil {
		log.Fatalf("Error when opening config file: %q", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	mattermostBot := bot.NewMattermostBot(configData)
	jiraHttpClient := jirahttpclient.NewJiraHttpClient(&http.Client{}, configData.JiraBaseUrl, configData.JiraBotUsername, configData.JiraBotPassword)
	jiraTaskCreator := jirataskcreator.NewJiraTaskCreator(jiraHttpClient)
	jiraCommentCreator := jiracommentcreator.NewJiraCommentCreator(jiraHttpClient)
	var provider service.MessagesProvider
	var sender service.MessagesSender
	if configData.Messenger == config.MATTERMOST {
		httpClientForMessanger := mattermosthttpclient.NewHttpClient(&http.Client{}, mattermostBot.Token, configData.MattermostHttp)
		provider = mattermostprovider.NewMattermostProvider(mattermostBot)
		sender = mattermostsender.NewMattermostSender(httpClientForMessanger)
	} else if configData.Messenger == config.ROCKETCHAT {
		client := ddp.NewClient("ws://"+configData.RocketchatHost+"/websocket", (&url.URL{Host: configData.RocketchatHost}).String())
		provider, err = rocketchatprovider.NewRocketChatPovider(client, configData.RocketchatId, configData.RocketchatToken)
		if err != nil {
			log.Fatalf("Rocket chat error: %q", err)
		}
		rocketchatHttpClient := rocketchathttpclient.NewHttpClient(&http.Client{}, configData.RocketchatId, configData.RocketchatToken, configData.RocketchatHttp)
		sender = rocketchatsender.NewRocketChatSender(rocketchatHttpClient)
	}
	matcher, err := messagesmatcher.NewMessagesMatcher(configData.MessagesPattern)
	if err != nil {
		log.Fatalf("Unsuccessful start: %q", err)
	}
	go provider.Run(ctx)
	var storage service.StorageMsgTasks
	if configData.EnableMsgThreating {
		if configData.StorageType == config.POSTGRES {
			storageValue, err := storagepostgres.NewStorageMsgTasksDB(ctx, configData.PostgresLogin,
				configData.PostgresPassword, configData.PostgresHost, configData.PostgresPort,
				configData.PostgresName, configData.PostgresTableName)
			if err != nil {
				log.Fatalf("Error when connect to postgres %q", err)
			}
			storage = &storageValue
		} else if configData.StorageType == config.IN_MEMORY {
			storageValue := storageinmemory.NewStorageMsgTasksInMemory()
			storage = &storageValue
		}
	}
	taskFromMessagesCreator := service.NewTaskFromMessagesCreator(provider, sender, matcher, jiraTaskCreator,
		configData.MessagesPatternTemplate, configData.JiraProject, configData.JiraIssueType,
		configData.EnableMsgThreating, storage, jiraCommentCreator)
	taskFromMessagesCreator.Run(ctx)
}
