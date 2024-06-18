package jiracommentcreator

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/jirahttpclient"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	tests := map[string]struct {
		text           string
		taskId         string
		username       string
		password       string
		expectedBase64 string
	}{
		"simple": {
			text:           "text",
			taskId:         "TEST-1",
			username:       "username",
			password:       "password",
			expectedBase64: "Basic dXNlcm5hbWU6cGFzc3dvcmQ=",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), "/rest/api/2/issue/"+tc.taskId+"/comment")
				bufSize := 1024
				buffer := make([]byte, bufSize)
				length, _ := req.Body.Read(buffer)
				assert.Equal(t, req.Header.Get("Authorization"), tc.expectedBase64)
				var requestComment jirahttpclient.JiraCommentRequest
				json.Unmarshal(buffer[:length], &requestComment)
				assert.Equal(t, requestComment.Body, tc.text)
				rw.Write([]byte(`OK`))
			}))
			defer server.Close()
			client := server.Client()
			jiraClient := jirahttpclient.NewJiraHttpClient(client, server.URL, tc.username, tc.password)
			commentCreator := NewJiraCommentCreator(jiraClient)
			commentCreator.CreateComment(context.Background(), tc.text, tc.taskId)
		})
	}
}
