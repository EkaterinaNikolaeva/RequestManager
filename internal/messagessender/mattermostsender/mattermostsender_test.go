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
	testMsgText := "test_post"
	testChannelId := "test_chat_id"
	testRootId := "01234"
	testAuthorId := "0"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/api/v4/posts")
		bufSize := 1024
		buffer := make([]byte, bufSize)
		length, _ := req.Body.Read(buffer)
		var post mattermosthttpclient.RequestPost
		json.Unmarshal(buffer[:length], &post)
		assert.Equal(t, post.ChannelId, testChannelId)
		assert.Equal(t, post.Message, testMsgText)
		assert.Equal(t, post.RootId, testRootId)
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	client := server.Client()
	mmClient := mattermosthttpclient.NewHttpClient(client, "", server.URL)
	sender := NewMattermostSender(mmClient)
	assert.Nil(t, sender.SendMessage(context.Background(), message.NewMessage(testMsgText, testChannelId, testRootId, message.NewMessageAuthor(testAuthorId))))
}

func TestSendMessagesWithError(t *testing.T) {
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
	assert.NotNil(t, sender.SendMessage(context.Background(), message.NewMessage("", "", "", message.NewMessageAuthor(""))))
}
