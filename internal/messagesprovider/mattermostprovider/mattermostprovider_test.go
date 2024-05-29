package mattermostprovider

import (
	"encoding/json"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
	"github.com/stretchr/testify/assert"
)

func TestGetMessage(t *testing.T) {
	tests := map[string]struct {
		msg            MattermostMessage
		wantErr        bool
		expectedResult message.Message
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
