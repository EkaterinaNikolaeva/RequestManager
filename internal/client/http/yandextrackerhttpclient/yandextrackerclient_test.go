package yandextrackerhttpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	task := RequestTask{
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
		var request RequestTask
		json.Unmarshal(buffer[:length], &request)
		assert.Equal(t, task, request)
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	client := server.Client()
	yandexTrackerHttpClient := NewYandexTracketHttpClient(server.URL, server.URL, "id-org", "TYPE-ORG", "Bearer", "token", client)
	yandexTrackerHttpClient.CreateTask(context.Background(), task)
}

func TestCreateComment(t *testing.T) {
	tests := map[string]struct {
		text             string
		idIssue          string
		idOrganization   string
		typeOrganization string
		tokenType        string
		token            string
		expectsError     bool
	}{
		"simple": {text: "text", idIssue: "Test-1", idOrganization: "id-org",
			typeOrganization: "TYPE-ORG", tokenType: "Bearer", token: "token", expectsError: false},
		"without header": {text: "other text", idIssue: "ID-3", idOrganization: "org",
			tokenType: "Oauth", token: "ttt", expectsError: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), "/v2/issues/"+tc.idIssue+"/comments")
				bufSize := 1024
				buffer := make([]byte, bufSize)
				length, _ := req.Body.Read(buffer)
				assert.Equal(t, req.Header.Get(tc.typeOrganization), tc.idOrganization)
				assert.Equal(t, req.Header.Get("Authorization"), tc.tokenType+" "+tc.token)
				var request RequestComment
				json.Unmarshal(buffer[:length], &request)
				assert.Equal(t, tc.text, request.Text)
				rw.Write([]byte(`OK`))
			}))
			defer server.Close()
			client := server.Client()
			yandexTrackerHttpClient := NewYandexTracketHttpClient(server.URL, server.URL, tc.idOrganization, tc.typeOrganization, tc.tokenType, tc.token, client)
			err := yandexTrackerHttpClient.AddComment(context.Background(), tc.text, tc.idIssue)
			if !tc.expectsError {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
