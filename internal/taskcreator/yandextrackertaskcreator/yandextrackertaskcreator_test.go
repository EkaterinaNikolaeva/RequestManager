package yandextrackertaskcreator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/yandextrackerhttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/task"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	task := task.TaskCreateRequest{
		Description: "Description",
		Type:        "Task",
		Name:        "Name",
		Project:     "Pr",
	}
	taskApiYandexTracker := yandextrackerhttpclient.RequestTask{
		Description: "Description",
		Type:        "Task",
		Summary:     "Name",
		Queue:       "Pr",
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/v2/issues/")
		bufSize := 1024
		buffer := make([]byte, bufSize)
		length, _ := req.Body.Read(buffer)
		assert.Equal(t, req.Header.Get("TYPE-ORG"), "id-org")
		assert.Equal(t, req.Header.Get("Authorization"), "Bearer token")
		var request yandextrackerhttpclient.RequestTask
		json.Unmarshal(buffer[:length], &request)
		assert.Equal(t, taskApiYandexTracker, request)
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	client := server.Client()
	yandexTrackerHttpClient := yandextrackerhttpclient.NewYandexTracketHttpClient(server.URL, server.URL, "id-org", "TYPE-ORG", "Bearer", "token", client)
	taskCreator := NewYandexTrackerTaskCreator(yandexTrackerHttpClient)
	taskCreator.CreateTask(task)
}
