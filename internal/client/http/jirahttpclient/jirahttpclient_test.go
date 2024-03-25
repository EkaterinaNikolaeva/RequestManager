package jirahttpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/jiratasks"
	"github.com/stretchr/testify/assert"
)

func TestSendMessages(t *testing.T) {
	summary := "Some summary"
	description := "More about issue"
	project := "PROJECT"
	issueType := "Bug"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/rest/api/2/issue/")
		bufSize := 1024
		buffer := make([]byte, bufSize)
		length, _ := req.Body.Read(buffer)
		var issue jiratasks.JiraTaskCreationRequest
		json.Unmarshal(buffer[:length], &issue)
		assert.Equal(t, issueType, issue.Fields.IssueType.Name)
		assert.Equal(t, summary, issue.Fields.Summary)
		assert.Equal(t, description, issue.Fields.Description)
		assert.Equal(t, project, issue.Fields.Project.Key)
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
	jiraClient.CreateTask(jiratasks.JiraTaskCreationRequest{
		Fields: jiratasks.JiraTaskCreationFields{
			Project: jiratasks.JiraTaskCreationProject{
				Key: project,
			},
			Summary:     summary,
			Description: description,
			IssueType: jiratasks.JiraTaskCreationIssueType{
				Name: issueType,
			},
		},
	})
}
