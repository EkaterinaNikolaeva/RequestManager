package rocketchatsender

import (
	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/rocketchathttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
)

type RocketChatSender struct {
	rocketChatHttpClient *rocketchathttpclient.RocketChatHttpClient
}

func NewRocketChatSender(httpClient *rocketchathttpclient.RocketChatHttpClient) RocketChatSender {
	return RocketChatSender{
		rocketChatHttpClient: httpClient,
	}
}

func (s RocketChatSender) SendMessage(message message.Message) error {
	return s.rocketChatHttpClient.SendMessage(mapRocketChatMessageFromMessage(message))
}

func mapRocketChatMessageFromMessage(message message.Message) rocketchathttpclient.RequestMessage {
	msg := rocketchathttpclient.RequestMessage{
		Message: rocketchathttpclient.RequestMessageData{
			Rid:  message.ChannelId,
			Tmid: message.RootMessageId,
			Msg:  message.MessageText,
		},
	}
	return msg
}
