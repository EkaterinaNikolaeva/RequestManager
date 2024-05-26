package mattermostsender

import (
	"context"

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

func (m MattermostSender) SendMessage(ctx context.Context, message message.Message) error {
	return m.mattermostHttpClient.CreatePost(ctx, mapMattermostPostFromMessage(message))
}

func mapMattermostPostFromMessage(message message.Message) mattermosthttpclient.RequestPost {
	post := mattermosthttpclient.RequestPost{
		Message:   message.MessageText,
		ChannelId: message.ChannelId,
		RootId:    message.RootMessageId,
	}
	return post
}
