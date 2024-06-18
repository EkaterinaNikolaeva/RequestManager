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
	tests := map[string]struct {
		username       string
		password       string
		expectedBase64 string
		resultId       string
		resultKey      string
		requestTask    JiraTaskCreationRequest
	}{
		"simple": {
			requestTask: JiraTaskCreationRequest{
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
			},
			resultId:       "000",
			resultKey:      "1",
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
				var issue JiraTaskCreationRequest
				json.Unmarshal(buffer[:length], &issue)

				assert.Equal(t, req.Header.Get("Authorization"), tc.expectedBase64)
				assert.Equal(t, tc.requestTask, issue)

				result := JiraTaskCreationResponse{
					Id:  tc.resultId,
					Key: tc.resultKey,
				}
				bytesResult, _ := json.Marshal(result)
				rw.Write(bytesResult)
			}))
			defer server.Close()
			client := server.Client()
			jiraClient := NewJiraHttpClient(client, server.URL, tc.username, tc.password)
			_, _, err := jiraClient.CreateTask(context.Background(), tc.requestTask)
			assert.Nil(t, err)
		})
	}
}

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
			comment := JiraCommentRequest{
				Body: "text",
			}
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), "/rest/api/2/issue/"+tc.taskId+"/comment")
				bufSize := 1024
				buffer := make([]byte, bufSize)
				length, _ := req.Body.Read(buffer)

				assert.Equal(t, req.Header.Get("Authorization"), tc.expectedBase64)
				var requestComment JiraCommentRequest
				json.Unmarshal(buffer[:length], &requestComment)
				assert.Equal(t, requestComment, comment)
				rw.Write([]byte(`OK`))
			}))
			defer server.Close()
			client := server.Client()
			jiraClient := NewJiraHttpClient(client, server.URL, tc.username, tc.password)
			jiraClient.AddComment(context.Background(), tc.text, tc.taskId)
		})
	}

}
