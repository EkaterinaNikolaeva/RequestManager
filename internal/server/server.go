package server

import (
	"log"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/mattermost/mattermost-server/v6/model"
)

func handler(event *model.WebSocketEvent) {
	if event.GetData()["post"] == nil {
		return
	}
	log.Println(event.GetData())
}

func MakeServer(mattermostBot bot.MattermostBot, url string) {
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
				handler(resp)
			}
		}
	}()
	select {}
}
