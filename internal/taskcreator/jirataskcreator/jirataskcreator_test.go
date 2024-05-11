package jirataskcreator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/api/jira/jiratasks"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/jirahttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/task"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	testTaskName := "test name"
	testDescription := "description"
	testProject := "TEST-PROJECT"
	testType := "bug"
	testTaskId := "0"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/rest/api/2/issue/")
		bufSize := 1024
		buffer := make([]byte, bufSize)
		length, _ := req.Body.Read(buffer)
		var task jiratasks.JiraTaskCreationRequest
		json.Unmarshal(buffer[:length], &task)
		assert.Equal(t, testTaskName, task.Fields.Summary)
		assert.Equal(t, testDescription, task.Fields.Description)
		assert.Equal(t, testProject, task.Fields.Project.Key)
		assert.Equal(t, testType, task.Fields.IssueType.Name)
		responseTask := jiratasks.JiraTaskCreationResponse{
			Id:  "000",
			Key: testProject + "-" + testTaskId,
		}
		bytes, err := json.Marshal(responseTask)
		assert.Nil(t, err)
		rw.Write(bytes)
	}))
	defer server.Close()
	client := server.Client()
	jiraClient := jirahttpclient.NewJiraHttpClient(client, server.URL, "", "")
	taskCreator := NewJiraTaskCreator(jiraClient)
	task, err := taskCreator.CreateTask(task.NewTaskCreateRequest(testTaskName, testDescription, testType, testProject))
	assert.Nil(t, err)
	assert.Equal(t, testDescription, task.Description)
	assert.Equal(t, testProject, task.Project)
	assert.Equal(t, testType, task.Type)
	assert.Equal(t, testTaskName, task.Name)
	assert.Equal(t, server.URL+"/projects/"+testProject+"/issues/"+testProject+"-"+testTaskId, task.Link)
}
