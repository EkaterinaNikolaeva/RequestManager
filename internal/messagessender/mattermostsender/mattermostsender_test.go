package mattermostsender

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/mattermosthttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
	"github.com/stretchr/testify/assert"
)

func TestSendMessages(t *testing.T) {
	tests := map[string]struct {
		msgText   string
		channelId string
		rootId    string
		authorId  string
		isBot     bool
	}{
		"correct": {
			msgText:   "test_post",
			channelId: "test_chat_id",
			rootId:    "test_chat_id",
			authorId:  "0",
		},
		"bot": {
			msgText:   "test_post",
			channelId: "test_chat_id",
			rootId:    "test_chat_id",
			authorId:  "0",
			isBot:     true,
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), "/api/v4/posts")
				bufSize := 1024
				buffer := make([]byte, bufSize)
				length, _ := req.Body.Read(buffer)
				var post mattermosthttpclient.RequestPost
				json.Unmarshal(buffer[:length], &post)

				assert.Equal(t, post.ChannelId, tc.channelId)
				assert.Equal(t, post.Message, tc.msgText)
				assert.Equal(t, post.RootId, tc.rootId)
				rw.Write([]byte(`OK`))
			}))
			defer server.Close()
			client := server.Client()
			mmClient := mattermosthttpclient.NewHttpClient(client, "", server.URL)
			sender := NewMattermostSender(mmClient)
			assert.Nil(t, sender.SendMessage(context.Background(),
				message.NewMessage(tc.msgText, tc.channelId, tc.rootId,
					message.NewMessageAuthor(tc.authorId, tc.isBot))))
		})
	}
}

func TestSendMessagesWithError(t *testing.T) {
	tests := map[string]struct {
		msgText   string
		channelId string
		rootId    string
		authorId  string
	}{
		"error": {
			msgText:   "test_post",
			channelId: "test_chat_id",
			rootId:    "test_chat_id",
			authorId:  "0",
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), "/api/v4/posts")
				bufSize := 1024
				buffer := make([]byte, bufSize)
				length, _ := req.Body.Read(buffer)
				var post mattermosthttpclient.RequestPost
				json.Unmarshal(buffer[:length], &post)
				rw.WriteHeader(404)
			}))
			defer server.Close()
			client := server.Client()
			mmClient := mattermosthttpclient.NewHttpClient(client, "", server.URL)
			sender := NewMattermostSender(mmClient)
			assert.NotNil(t, sender.SendMessage(context.Background(),
				message.NewMessage(tc.msgText, tc.channelId, tc.rootId,
					message.NewMessageAuthor(tc.authorId))))
		})
	}
}
