package rocketchatsender

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/rocketchathttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
	"github.com/stretchr/testify/assert"
)

func TestSendMessages(t *testing.T) {
	tests := map[string]struct {
		msgText   string
		channelId string
		rootId    string
		authorId  string
		token     string
		userId    string
	}{
		"correct": {
			msgText:   "test_post",
			channelId: "test_chat_id",
			rootId:    "test_chat_id",
			authorId:  "0",
			token:     "token",
			userId:    "user-id",
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				msgRocketChat := rocketchathttpclient.RequestMessage{
					Message: rocketchathttpclient.RequestMessageData{
						Msg:  tc.msgText,
						Tmid: tc.rootId,
						Rid:  tc.channelId,
					},
				}
				assert.Equal(t, req.URL.String(), "/api/v1/chat.sendMessage")
				bufSize := 1024
				buffer := make([]byte, bufSize)
				length, _ := req.Body.Read(buffer)

				assert.Equal(t, req.Header.Get("X-User-Id"), tc.userId)
				assert.Equal(t, req.Header.Get("X-Auth-Token"), tc.token)
				var request rocketchathttpclient.RequestMessage
				json.Unmarshal(buffer[:length], &request)
				assert.Equal(t, msgRocketChat, request)
				rw.Write([]byte(`OK`))
			}))
			defer server.Close()
			client := server.Client()
			rocketChatClient := rocketchathttpclient.NewHttpClient(client, tc.userId, tc.token, server.URL)
			sender := NewRocketChatSender(rocketChatClient)
			assert.Nil(t, sender.SendMessage(context.Background(),
				message.NewMessage(tc.msgText, tc.channelId, tc.rootId,
					message.NewMessageAuthor(tc.authorId))))
		})
	}
}
