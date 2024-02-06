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

func MakeServer(mattermostBot bot.MattermostBot) {
	httpClient := mattermostmessages.NewHttpClient(&http.Client{})
	mattermostClient := model.NewAPIv4Client(mattermostBot.MattermostHttp)
	mattermostClient.SetToken(mattermostBot.Token)
	mattermostClient.MockSession(mattermostBot.Token)
	log.Println(mattermostClient.GetTeamByName(mattermostBot.TeamName, ""))
	webSocketClient, err := model.NewWebSocketClient4(mattermostBot.MattermostWebsocket, mattermostBot.Token)
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
