package yandextrackerhttpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apiyandextracker "github.com/EkaterinaNikolaeva/RequestManager/internal/api/yandextracker"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	task := apiyandextracker.RequestTask{
		Queue:       "TEST-QUEUE",
		Summary:     "Summary",
		Description: "Description",
		Type:        "Task",
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/v2/issues/")
		bufSize := 1024
		buffer := make([]byte, bufSize)
		length, _ := req.Body.Read(buffer)
		assert.Equal(t, req.Header.Get("TYPE-ORG"), "id-org")
		assert.Equal(t, req.Header.Get("Authorization"), "Bearer token")
		var request apiyandextracker.RequestTask
		json.Unmarshal(buffer[:length], &request)
		assert.Equal(t, task, request)
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	client := server.Client()
	yandexTrackerHttpClient := NewYandexTracketHttpClient(server.URL, server.URL, "id-org", "TYPE-ORG", "Bearer", "token", client)
	yandexTrackerHttpClient.CreateTask(task)
}

func TestCreateComment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/v2/issues/TEST-1/comments")
		bufSize := 1024
		buffer := make([]byte, bufSize)
		length, _ := req.Body.Read(buffer)
		assert.Equal(t, req.Header.Get("TYPE-ORG"), "id-org")
		assert.Equal(t, req.Header.Get("Authorization"), "Bearer token")
		var request apiyandextracker.RequestComment
		json.Unmarshal(buffer[:length], &request)
		assert.Equal(t, "text", request.Text)
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	client := server.Client()
	yandexTrackerHttpClient := NewYandexTracketHttpClient(server.URL, server.URL, "id-org", "TYPE-ORG", "Bearer", "token", client)
	yandexTrackerHttpClient.AddComment("text", "TEST-1")
}
