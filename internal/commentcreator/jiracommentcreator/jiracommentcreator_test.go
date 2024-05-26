package jiracommentcreator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/jirahttpclient"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	comment := jirahttpclient.JiraCommentRequest{
		Body: "text",
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/rest/api/2/issue/"+"TEST-1"+"/comment")
		bufSize := 1024
		buffer := make([]byte, bufSize)
		length, _ := req.Body.Read(buffer)
		assert.Equal(t, req.Header.Get("Authorization"), "Basic dXNlcm5hbWU6cGFzc3dvcmQ=")
		var requestComment jirahttpclient.JiraCommentRequest
		json.Unmarshal(buffer[:length], &requestComment)
		assert.Equal(t, requestComment, comment)
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	client := server.Client()
	jiraClient := jirahttpclient.NewJiraHttpClient(client, server.URL, "username", "password")
	commentCreator := NewJiraCommentCreator(jiraClient)
	commentCreator.CreateComment("text", "TEST-1")
}
