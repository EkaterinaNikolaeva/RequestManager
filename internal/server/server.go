package server

import (
	"log"
	"net/http"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/mattermostmessages"
	"github.com/mattermost/mattermost-server/v6/model"
)

func checkEventForJiraRequest(event *model.WebSocketEvent, client *mattermostmessages.HttpClient, mattermostBot bot.MattermostBot) {
	if event.GetData()["post"] == nil {
		return
	}
	client.CheckMessageForJiraRequest(event.GetData()["post"].(string), mattermostBot)
}

func MakeServer(mattermostBot bot.MattermostBot, url string) {
	httpClient := mattermostmessages.NewHttpClient(&http.Client{})
	mattermostClient := model.NewAPIv4Client("http://localhost:8065")
	mattermostClient.SetToken(mattermostBot.Token)
	mattermostClient.MockSession(mattermostBot.Token)
	log.Println(mattermostClient.GetTeamByName("jira-mattermost", ""))
	webSocketClient, err := model.NewWebSocketClient4("ws://localhost:8065", mattermostBot.Token)
	if err != nil {
		log.Printf("error when creating new websocket client: %q", err)
	}
	webSocketClient.Listen()
	go func() {
		for {
			select {
			case resp := <-webSocketClient.EventChannel:
				checkEventForJiraRequest(resp, httpClient, mattermostBot)
			}
		}
	}()
	select {}
}
