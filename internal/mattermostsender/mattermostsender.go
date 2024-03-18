package mattermostsender

import (
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/mattermostmessages"
)

type MattermostSender struct {
	mattermostHttpClient *mattermostmessages.MattermostHttpClient
}

func NewMattermostSender(httpClient *mattermostmessages.MattermostHttpClient) MattermostSender {
	return MattermostSender{
		mattermostHttpClient: httpClient,
	}
}

func (m MattermostSender) SendMessage(message message.Message) error {
	return m.mattermostHttpClient.CreatePost(mapMattermostPostFromMessage(message))
}

func mapMattermostPostFromMessage(message message.Message) mattermostmessages.RequestPost {
	post := mattermostmessages.RequestPost{
		Message:   message.MessageText,
		ChannelId: message.ChannelId,
		RootId:    message.RootMessageId,
	}
	return post
}
