package mattermostprovider

import (
	"encoding/json"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
)

type MattermostMessage struct {
	Id            string                 `json:"id,omitempty"`
	CreateAt      int                    `json:"create_at,omitempty"`
	UpdateAt      int                    `json:"update_at,omitempty"`
	DeleteAt      int                    `json:"delete_at,omitempty"`
	EditAt        int                    `json:"edit_at,omitempty"`
	UserId        string                 `json:"user_id,omitempty"`
	ChannelId     string                 `json:"channel_id,omitempty"`
	RootId        string                 `json:"root_id,omitempty"`
	OriginalId    string                 `json:"original_id,omitempty"`
	Message       string                 `json:"message,omitempty"`
	Type          string                 `json:"type,omitempty"`
	Props         map[string]interface{} `json:"props,omitempty"`
	Hashtag       string                 `json:"hashtag,omitempty"`
	PendingPostId string                 `json:"pending_post_id,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

func getFrom(msg MattermostMessage) string {
	rootId := msg.Id
	if msg.RootId != "" {
		rootId = msg.RootId
	}
	return rootId
}

func checkMessageFromBot(msg MattermostMessage) bool {
	props := msg.Props
	fromBot, ok := props["from_bot"]
	isBot, isBool := fromBot.(bool)
	return ok && isBool && isBot || fromBot == "true"
}

func GetMessage(bytes string) (message.Message, error) {
	var msg MattermostMessage
	err := json.Unmarshal([]byte(bytes), &msg)
	if err != nil {
		return message.Message{}, err
	}
	return message.Message{
		MessageText:   msg.Message,
		ChannelId:     msg.ChannelId,
		RootMessageId: getFrom(msg),
		Author:        message.MessageAuthor{Id: msg.UserId, IsBot: checkMessageFromBot(msg)},
	}, nil
}
