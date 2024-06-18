package mattermosthttpclient

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
		text      string
		channelId string
		tmid      string
		userId    string
		token     string
	}{
		"simple": {
			text:      "text",
			channelId: "chat-id",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			requestPost := RequestPost{
				Message:   tc.text,
				ChannelId: tc.channelId,
			}
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), "/api/v4/posts")
				bufSize := 1024
				buffer := make([]byte, bufSize)
				length, _ := req.Body.Read(buffer)
				var post RequestPost
				json.Unmarshal(buffer[:length], &post)

				assert.Equal(t, requestPost, post)
				rw.Write([]byte(`OK`))
			}))
			defer server.Close()
			client := server.Client()
			assert.Nil(t, NewHttpClient(client, "", server.URL).CreatePost(context.Background(), requestPost))
		})
	}

}
