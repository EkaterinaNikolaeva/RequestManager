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
	channel       chan service.Message
	httpClient    *mattermostmessages.HttpClient
	mattermostBot bot.MattermostBot
}

func NewMattermostProvider(mattermostBot bot.MattermostBot) MattermostProvider {
	var provider MattermostProvider
	provider.channel = make(chan service.Message)
	provider.httpClient = mattermostmessages.NewHttpClient(&http.Client{})
	provider.mattermostBot = mattermostBot
	return provider
}

func (m MattermostProvider) getMessageHandler(event *model.WebSocketEvent) {
	if event.GetData()["post"] == nil {
		return
	}
	message, from_bot, err := m.httpClient.GetMessage(event.GetData()["post"].(string))
	if err != nil {
		log.Printf("error when encoding message: %q", err)
	}
	if !from_bot {
		m.channel <- message
	}
}

func (m MattermostProvider) SendMessage(message service.Message) error {
	return m.httpClient.SendMessage(message, m.mattermostBot)
}

func (m MattermostProvider) Run() {
	mattermostClient := model.NewAPIv4Client(m.mattermostBot.MattermostHttp)
	mattermostClient.SetToken(m.mattermostBot.Token)
	mattermostClient.MockSession(m.mattermostBot.Token)
	webSocketClient, err := model.NewWebSocketClient4(m.mattermostBot.MattermostWebsocket, m.mattermostBot.Token)
	if err != nil {
		log.Printf("error when creating new websocket client: %q", err)
	}
	webSocketClient.Listen()
	for {
		select {
		case resp := <-webSocketClient.EventChannel:
			m.getMessageHandler(resp)
		}
	}

}

func (m MattermostProvider) GetMessagesChannel() <-chan service.Message {
	return m.channel
}
