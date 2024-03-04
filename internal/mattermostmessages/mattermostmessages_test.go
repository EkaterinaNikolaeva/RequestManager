package mattermostmessages

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessages(t *testing.T) {
	testMsg := "test_post"
	testChannelId := "test_chat_id"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/api/v4/posts")
		bufSize := 1024
		buffer := make([]byte, bufSize)
		length, _ := req.Body.Read(buffer)
		var post RequestPost
		json.Unmarshal(buffer[:length], &post)
		assert.Equal(t, post.ChannelId, testChannelId)
		assert.Equal(t, post.Message, testMsg)
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	client := server.Client()
	NewHttpClient(client, server.URL, "").CreatePost(RequestPost{
		Message:   testMsg,
		ChannelId: testChannelId,
	})
}
