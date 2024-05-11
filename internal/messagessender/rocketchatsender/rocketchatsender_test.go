package rocketchatsender

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apirocketchat "github.com/EkaterinaNikolaeva/RequestManager/internal/api/rocketchat"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/rocketchathttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
	"github.com/stretchr/testify/assert"
)

func TestSendMessages(t *testing.T) {
	msg := message.Message{
		MessageText:   "text",
		Author:        message.MessageAuthor{Id: "id-user", IsBot: true},
		ChannelId:     "channel",
		RootMessageId: "root-msg",
	}
	msgRocketChat := apirocketchat.RequestMessage{
		Message: apirocketchat.RequestMessageData{
			Msg:  "text",
			Tmid: "root-msg",
			Rid:  "channel",
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/api/v1/chat.sendMessage")
		bufSize := 1024
		buffer := make([]byte, bufSize)
		length, _ := req.Body.Read(buffer)
		assert.Equal(t, req.Header.Get("X-User-Id"), "user-id")
		assert.Equal(t, req.Header.Get("X-Auth-Token"), "token")
		var request apirocketchat.RequestMessage
		json.Unmarshal(buffer[:length], &request)
		assert.Equal(t, msgRocketChat, request)
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	client := server.Client()
	rocketChatClient := rocketchathttpclient.NewHttpClient(client, "user-id", "token", server.URL)
	sender := NewRocketChatSender(rocketChatClient)
	assert.Nil(t, sender.SendMessage(msg))
}
