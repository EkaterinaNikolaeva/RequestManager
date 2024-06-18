package rocketchathttpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessages(t *testing.T) {
	tests := map[string]struct {
		text   string
		rid    string
		tmid   string
		userId string
		token  string
	}{
		"simple": {
			text:   "text",
			rid:    "chat-id",
			tmid:   "root-msg-id",
			userId: "user-id",
			token:  "token",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg := RequestMessage{
				Message: RequestMessageData{Msg: tc.text,
					Rid:  tc.rid,
					Tmid: tc.tmid,
				},
			}
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), "/api/v1/chat.sendMessage")
				bufSize := 1024
				buffer := make([]byte, bufSize)
				length, _ := req.Body.Read(buffer)
				assert.Equal(t, req.Header.Get("X-User-Id"), tc.userId)
				assert.Equal(t, req.Header.Get("X-Auth-Token"), tc.token)
				var request RequestMessage
				json.Unmarshal(buffer[:length], &request)
				assert.Equal(t, msg, request)
				rw.Write([]byte(`OK`))
			}))
			defer server.Close()
			client := server.Client()
			rocketChatClient := NewHttpClient(client, tc.userId, tc.token, server.URL)
			assert.Nil(t, rocketChatClient.SendMessage(context.Background(), msg))
		})
	}

}
