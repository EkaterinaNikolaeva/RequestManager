package mattermostsender

import (
	"github.com/EkaterinaNikolaeva/RequestManager/internal/mattermostmessages"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/message"
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
	return m.mattermostHttpClient.SendMessage(message)
}
