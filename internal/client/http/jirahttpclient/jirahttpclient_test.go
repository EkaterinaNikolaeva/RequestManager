package jirahttpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/jiratasks"
	"github.com/stretchr/testify/assert"
)

func TestCreateIssue(t *testing.T) {
	requestTask := jiratasks.JiraTaskCreationRequest{
		Fields: jiratasks.JiraTaskCreationFields{
			Project: jiratasks.JiraTaskCreationProject{
				Key: "PROJECT",
			},
			Summary:     "Some summary",
			Description: "More about issue",
			IssueType: jiratasks.JiraTaskCreationIssueType{
				Name: "Bug",
			},
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/rest/api/2/issue/")
		bufSize := 1024
		buffer := make([]byte, bufSize)
		length, _ := req.Body.Read(buffer)
		var issue jiratasks.JiraTaskCreationRequest
		json.Unmarshal(buffer[:length], &issue)
		assert.Equal(t, requestTask, issue)
		result := jiratasks.JiraTaskCreationResponse{
			Id:  "000",
			Key: "1",
		}
		bytesResult, _ := json.Marshal(result)
		rw.Write(bytesResult)
	}))
	defer server.Close()
	client := server.Client()
	jiraClient := NewJiraHttpClient(client, server.URL, "", "")
	_, _, err := jiraClient.CreateTask(requestTask)
	assert.Nil(t, err)
}