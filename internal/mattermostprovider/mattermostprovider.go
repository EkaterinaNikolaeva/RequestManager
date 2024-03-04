package mattermostprovider

import (
	"context"
	"log"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/mattermostmessages"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/message"
	"github.com/mattermost/mattermost-server/v6/model"
)

type MattermostProvider struct {
	channel       chan message.Message
	mattermostBot bot.MattermostBot
}

func NewMattermostProvider(mattermostBot bot.MattermostBot) MattermostProvider {
	var provider MattermostProvider
	provider.channel = make(chan message.Message)
	provider.mattermostBot = mattermostBot
	return provider
}

func checkRequestPostType(event *model.WebSocketEvent) bool {
	return event.GetData()["post"] != nil
}

func (m MattermostProvider) handleMessage(event *model.WebSocketEvent) {
	if !checkRequestPostType(event) {
		return
	}
	message, err := mattermostmessages.GetMessage(event.GetData()["post"].(string))
	if err != nil {
		log.Printf("error when encoding message: %q", err)
		return
	}
	m.channel <- message
}

func (m MattermostProvider) Run(ctx context.Context) {
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
		case <-ctx.Done():
			log.Printf("ctx is done, stop websocket mattermost client")
			webSocketClient.Close()
			return
		case resp := <-webSocketClient.EventChannel:
			m.handleMessage(resp)
		}
	}
}

func (m MattermostProvider) GetMessagesChannel() <-chan message.Message {
	return m.channel
}
