package jirataskcreator

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/jirahttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/task"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	tests := map[string]struct {
		username            string
		password            string
		expectedBase64      string
		resultId            string
		resultKey           string
		expectedRequestTask jirahttpclient.JiraTaskCreationRequest
	}{
		"simple": {
			expectedRequestTask: jirahttpclient.JiraTaskCreationRequest{
				Fields: jirahttpclient.JiraTaskCreationFields{
					Project: jirahttpclient.JiraTaskCreationProject{
						Key: "PROJECT",
					},
					Summary:     "Some summary",
					Description: "More about issue",
					IssueType: jirahttpclient.JiraTaskCreationIssueType{
						Name: "Bug",
					},
				},
			},
			resultId:       "000",
			username:       "username",
			password:       "password",
			expectedBase64: "Basic dXNlcm5hbWU6cGFzc3dvcmQ=",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), "/rest/api/2/issue/")
				bufSize := 1024
				buffer := make([]byte, bufSize)
				length, _ := req.Body.Read(buffer)
				var task jirahttpclient.JiraTaskCreationRequest
				json.Unmarshal(buffer[:length], &task)
				assert.Equal(t, tc.expectedRequestTask, task)
				responseTask := jirahttpclient.JiraTaskCreationResponse{
					Id:  tc.resultId,
					Key: tc.expectedRequestTask.Fields.Project.Key + "-" + tc.resultId,
				}
				bytes, err := json.Marshal(responseTask)
				assert.Nil(t, err)
				rw.Write(bytes)
			}))
			defer server.Close()
			client := server.Client()
			jiraClient := jirahttpclient.NewJiraHttpClient(client, server.URL, "", "")
			taskCreator := NewJiraTaskCreator(jiraClient)
			res, err := taskCreator.CreateTask(context.Background(), task.NewTaskCreateRequest(tc.expectedRequestTask.Fields.Summary, tc.expectedRequestTask.Fields.Description,
				tc.expectedRequestTask.Fields.IssueType.Name,
				tc.expectedRequestTask.Fields.Project.Key))
			assert.Nil(t, err)
			assert.Equal(t, task.TaskCreated{
				Name:        tc.expectedRequestTask.Fields.Summary,
				Description: tc.expectedRequestTask.Fields.Description,
				Type:        tc.expectedRequestTask.Fields.IssueType.Name,
				Project:     tc.expectedRequestTask.Fields.Project.Key,
				Link:        server.URL + "/projects/" + tc.expectedRequestTask.Fields.Project.Key + "/issues/" + tc.expectedRequestTask.Fields.Project.Key + "-" + tc.resultId,
				Id:          tc.expectedRequestTask.Fields.Project.Key + "-" + tc.resultId,
			}, res)
		})
	}
}
