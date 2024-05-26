package rocketchathttpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessages(t *testing.T) {
	msg := RequestMessage{
		Message: RequestMessageData{Msg: "text",
			Rid:  "chat-id",
			Tmid: "root-msg-id",
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/api/v1/chat.sendMessage")
		bufSize := 1024
		buffer := make([]byte, bufSize)
		length, _ := req.Body.Read(buffer)
		assert.Equal(t, req.Header.Get("X-User-Id"), "user-id")
		assert.Equal(t, req.Header.Get("X-Auth-Token"), "token")
		var request RequestMessage
		json.Unmarshal(buffer[:length], &request)
		assert.Equal(t, msg, request)
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	client := server.Client()
	rocketChatClient := NewHttpClient(client, "user-id", "token", server.URL)
	assert.Nil(t, rocketChatClient.SendMessage(msg))
}
