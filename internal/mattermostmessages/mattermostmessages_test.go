package mattermostmessages

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/stretchr/testify/assert"
)

func TestSendMessages(t *testing.T) {
	testMsg := "test_msg"
	testChannelId := "test_chat_id"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/api/v4/posts")
		bufSize := 1024
		buffer := make([]byte, bufSize)
		length, _ := req.Body.Read(buffer)
		var msg Message
		json.Unmarshal(buffer[:length], &msg)
		assert.Equal(t, msg.ChannelId, testChannelId)
		assert.Equal(t, msg.Message, testMsg)
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	client := server.Client()
	mattermostBot := bot.LoadMattermostBot()
	NewHttpClient(client).SendMessage(Message{
		Message:   testMsg,
		ChannelId: testChannelId,
	}, server.URL, mattermostBot)
}
