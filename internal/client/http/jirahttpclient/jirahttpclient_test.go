package jirahttpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIssue(t *testing.T) {
	requestTask := JiraTaskCreationRequest{
		Fields: JiraTaskCreationFields{
			Project: JiraTaskCreationProject{
				Key: "PROJECT",
			},
			Summary:     "Some summary",
			Description: "More about issue",
			IssueType: JiraTaskCreationIssueType{
				Name: "Bug",
			},
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/rest/api/2/issue/")
		bufSize := 1024
		buffer := make([]byte, bufSize)
		length, _ := req.Body.Read(buffer)
		var issue JiraTaskCreationRequest
		json.Unmarshal(buffer[:length], &issue)
		assert.Equal(t, requestTask, issue)
		result := JiraTaskCreationResponse{
			Id:  "000",
			Key: "1",
		}
		bytesResult, _ := json.Marshal(result)
		rw.Write(bytesResult)
	}))
	defer server.Close()
	client := server.Client()
	jiraClient := NewJiraHttpClient(client, server.URL, "", "")
	_, _, err := jiraClient.CreateTask(context.Background(), requestTask)
	assert.Nil(t, err)
}

func TestCreateComment(t *testing.T) {
	comment := JiraCommentRequest{
		Body: "text",
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/rest/api/2/issue/"+"TEST-1"+"/comment")
		bufSize := 1024
		buffer := make([]byte, bufSize)
		length, _ := req.Body.Read(buffer)
		assert.Equal(t, req.Header.Get("Authorization"), "Basic dXNlcm5hbWU6cGFzc3dvcmQ=")
		var requestComment JiraCommentRequest
		json.Unmarshal(buffer[:length], &requestComment)
		assert.Equal(t, requestComment, comment)
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	client := server.Client()
	jiraClient := NewJiraHttpClient(client, server.URL, "username", "password")
	jiraClient.AddComment(context.Background(), "text", "TEST-1")
}
