package mattermostsender

import (
	"github.com/EkaterinaNikolaeva/RequestManager/internal/api/mattermost/mattermostmessages"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/mattermosthttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
)

type MattermostSender struct {
	mattermostHttpClient *mattermosthttpclient.MattermostHttpClient
}

func NewMattermostSender(httpClient *mattermosthttpclient.MattermostHttpClient) MattermostSender {
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
