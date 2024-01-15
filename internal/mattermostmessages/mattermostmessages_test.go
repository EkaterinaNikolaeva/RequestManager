package mattermostmessages

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/stretchr/testify/assert"
)

func TestSendMessages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/api/v4/posts")
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	client := server.Client()
	mattermostBot := bot.LoadMattermostBot()
	NewHttpClient(client).SendMessage(Message{
		Message:   "test_msg",
		ChannelId: "test_chat_id",
	}, server.URL, mattermostBot)
}
