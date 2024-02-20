package mattermostprovider

import (
	"log"
	"net/http"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/mattermostmessages"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/service"
	"github.com/mattermost/mattermost-server/v6/model"
)

type MattermostProvider struct {
	channel chan service.Message
}

func NewMattermostProvider() MattermostProvider {
	var provider MattermostProvider
	provider.channel = make(chan service.Message)
	return provider
}

func (m MattermostProvider) getMessageHandler(event *model.WebSocketEvent, client *mattermostmessages.HttpClient, mattermostBot bot.MattermostBot) {
	if event.GetData()["post"] == nil {
		return
	}
	message, err := client.GetMessage(event.GetData()["post"].(string))
	if err != nil {
		log.Printf("error when encoding message: %q", err)
	}
	m.channel <- service.Message{Message: message}
}

func (m MattermostProvider) Run(mattermostBot bot.MattermostBot) {
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
				m.getMessageHandler(resp, httpClient, mattermostBot)
			}
		}
	}()
	select {}
}

func (m MattermostProvider) GetMessagesChannel() <-chan service.Message {
	return m.channel
}
