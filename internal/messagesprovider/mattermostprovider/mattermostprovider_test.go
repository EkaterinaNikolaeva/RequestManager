package mattermostprovider

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/config"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/stretchr/testify/assert"
)

func TestGetMessage(t *testing.T) {
	tests := map[string]struct {
		msg            MattermostMessage
		wantErr        bool
		expectedResult message.Message
		websocket      message.Message
	}{
		"simple": {
			msg: MattermostMessage{
				Id:         "id",
				CreateAt:   111,
				UpdateAt:   111,
				DeleteAt:   0,
				EditAt:     0,
				UserId:     "user-id",
				ChannelId:  "channel",
				RootId:     "root",
				OriginalId: "original",
				Message:    "text",
				Type:       "msg",
			},
			wantErr: false,
			expectedResult: message.Message{
				MessageText:   "text",
				ChannelId:     "channel",
				RootMessageId: "root",
				Author: message.MessageAuthor{
					Id:    "user-id",
					IsBot: false,
				},
			},
		},
		"bot": {
			msg: MattermostMessage{
				Id:         "id",
				CreateAt:   111,
				UpdateAt:   111,
				DeleteAt:   0,
				EditAt:     0,
				UserId:     "user-id",
				ChannelId:  "channel",
				RootId:     "root",
				OriginalId: "original",
				Message:    "text",
				Type:       "msg",
				Props: map[string]interface{}{
					"is_bot": true,
				},
			},
			wantErr: false,
			expectedResult: message.Message{
				MessageText:   "text",
				ChannelId:     "channel",
				RootMessageId: "root",
				Author: message.MessageAuthor{
					Id:    "user-id",
					IsBot: false,
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bytes, err := json.Marshal(tc.msg)
			assert.Nil(t, err)
			domainMessage, err := GetMessage(string(bytes))
			if (err != nil) != tc.wantErr {
				t.Errorf("TestGetMessage() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			assert.Equal(t, tc.expectedResult, domainMessage)
		})
	}
}

func TestHandleMessage(t *testing.T) {
	mmProvider := NewMattermostProvider(bot.NewMattermostBot(config.Config{}))
	tests := map[string]struct {
		msg            MattermostMessage
		wantErr        bool
		expectedResult message.Message
		websocket      message.Message
		mmProvider     MattermostProvider
	}{
		"simple": {
			msg: MattermostMessage{
				Id:         "id",
				CreateAt:   111,
				UpdateAt:   111,
				DeleteAt:   0,
				EditAt:     0,
				UserId:     "user-id",
				ChannelId:  "channel",
				RootId:     "root",
				OriginalId: "original",
				Message:    "text",
				Type:       "msg",
			},
			wantErr: false,
			expectedResult: message.Message{
				MessageText:   "text",
				ChannelId:     "channel",
				RootMessageId: "root",
				Author: message.MessageAuthor{
					Id:    "user-id",
					IsBot: false,
				},
			},
			mmProvider: mmProvider,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bytes, err := json.Marshal(tc.msg)
			assert.Nil(t, err)
			domainMessage, err := GetMessage(string(bytes))
			if (err != nil) != tc.wantErr {
				t.Errorf("TestGetMessage() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			assert.Equal(t, tc.expectedResult, domainMessage)
			ws := model.NewWebSocketEvent("event", "team", "channel", "user-id", map[string]bool{})
			bytes2, _ := json.Marshal(tc.msg)
			ws.Add("post", string(bytes2))
			go tc.mmProvider.handleMessage(ws)
			actualMessage := ReadFromChannelWithTimeout[message.Message](t, mmProvider.channel, time.Second)
			assert.Equal(t, tc.expectedResult, actualMessage)
		})
	}
}

func ReadFromChannelWithTimeout[T any](t *testing.T, channel <-chan T, timeout time.Duration) T {
	select {
	case msg := <-channel:
		return msg
	case <-time.NewTimer(timeout).C:
		t.Error("Expected message not received")

		var t T

		return t
	}

}
