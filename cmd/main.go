package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/commentcreator/jiracommentcreator"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/commentcreator/yandextrackercommentcreator"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/config"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/messagesmatcher"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/messagesprovider/mattermostprovider"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/messagesprovider/rocketchatprovider"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/messagessender/mattermostsender"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/messagessender/rocketchatsender"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/service"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/storage/storageinmemory"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/storage/storagepostgres"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/taskcreator/jirataskcreator"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/taskcreator/yandextrackertaskcreator"
	"github.com/gopackage/ddp"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/jirahttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/mattermosthttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/rocketchathttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/yandextrackerhttpclient"

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
	var taskCreator service.TaskCreator
	var commentCreator service.CommentCreator
	var defaultProject string
	var defaultTypeTask string

	if configData.TaskTracker == config.TaskTrackerJira {
		jiraHttpClient := jirahttpclient.NewJiraHttpClient(&http.Client{}, configData.JiraBaseUrl, configData.JiraBotUsername, configData.JiraBotPassword)
		taskCreator = jirataskcreator.NewJiraTaskCreator(jiraHttpClient)
		commentCreator = jiracommentcreator.NewJiraCommentCreator(jiraHttpClient)
		defaultProject = configData.JiraProject
		defaultTypeTask = configData.JiraIssueType
	} else if configData.TaskTracker == config.TaskTrackerYandexTracker {
		yandexTrackerHttpClient := yandextrackerhttpclient.NewYandexTracketHttpClient(configData.YandexTrackerHost,
			configData.YandexTrackerBaseUrl, configData.YandexTrackerIdOrganization, configData.YandexTrackerTypeOrganization,
			configData.YandexTrackerTokenType, configData.YandexTrackerToken, &http.Client{})
		taskCreator = yandextrackertaskcreator.NewYandexTrackerTaskCreator(yandexTrackerHttpClient)
		commentCreator = yandextrackercommentcreator.NewYandexTrackerCommentCreator(yandexTrackerHttpClient)
		defaultProject = configData.YandexTrackerQueue
		defaultTypeTask = configData.YandexTrackerTaskType
	}
	var provider service.MessagesProvider
	var sender service.MessagesSender
	if configData.Messenger == config.MessengerMattermost {
		mattermostBot := bot.NewMattermostBot(configData)
		httpClientForMessanger := mattermosthttpclient.NewHttpClient(&http.Client{}, mattermostBot.Token, configData.MattermostHttp)
		provider = mattermostprovider.NewMattermostProvider(mattermostBot)
		sender = mattermostsender.NewMattermostSender(httpClientForMessanger)
	} else if configData.Messenger == config.MessengerRocketChat {
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
	taskFromMessagesCreator := service.NewTaskFromMessagesCreator(provider, sender, matcher, taskCreator,
		configData.MessagesPatternTemplate, defaultProject, defaultTypeTask,
		configData.EnableMsgThreating, storage, commentCreator, configData.Messenger, configData.TaskNamePatternTemplate)
	taskFromMessagesCreator.Run(ctx)
}
