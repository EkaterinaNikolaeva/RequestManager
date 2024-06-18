package yandextrackertaskcreator

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/yandextrackerhttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/task"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	tests := map[string]struct {
		idOrganization    string
		typeOrganization  string
		tokenType         string
		token             string
		taskCreateRequest task.TaskCreateRequest
	}{
		"simple": {
			idOrganization:   "id-org",
			typeOrganization: "TYPE-ORG",
			tokenType:        "Bearer",
			token:            "token",

			taskCreateRequest: task.TaskCreateRequest{
				Description: "Description",
				Type:        "Task",
				Name:        "Name",
				Project:     "Pr",
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			taskApiYandexTracker := yandextrackerhttpclient.RequestTask{
				Description: tc.taskCreateRequest.Description,
				Type:        tc.taskCreateRequest.Type,
				Summary:     tc.taskCreateRequest.Name,
				Queue:       tc.taskCreateRequest.Project,
			}
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), "/v2/issues/")
				bufSize := 1024
				buffer := make([]byte, bufSize)
				length, _ := req.Body.Read(buffer)
				assert.Equal(t, req.Header.Get(tc.typeOrganization), tc.idOrganization)
				assert.Equal(t, req.Header.Get("Authorization"), tc.tokenType+" "+tc.token)
				var request yandextrackerhttpclient.RequestTask
				json.Unmarshal(buffer[:length], &request)
				assert.Equal(t, taskApiYandexTracker, request)
				rw.Write([]byte(`OK`))
			}))
			defer server.Close()
			client := server.Client()
			yandexTrackerHttpClient := yandextrackerhttpclient.NewYandexTracketHttpClient(server.URL,
				server.URL,
				tc.idOrganization,
				tc.typeOrganization,
				tc.tokenType,
				tc.token,
				client)
			taskCreator := NewYandexTrackerTaskCreator(yandexTrackerHttpClient)
			taskCreator.CreateTask(context.Background(), tc.taskCreateRequest)
		})
	}

}
