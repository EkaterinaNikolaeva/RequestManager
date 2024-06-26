package yandextrackercommentcreator

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/yandextrackerhttpclient"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	tests := map[string]struct {
		typeOrg           string
		idOrg             string
		typeAuthorization string
		token             string
		text              string
		taskId            string
	}{
		"simple": {
			typeOrg:           "TYPE-ORG",
			idOrg:             "id-org",
			typeAuthorization: "Bearer",
			token:             "token",
			text:              "text",
			taskId:            "TEST-1",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), "/v2/issues/"+tc.taskId+"/comments")
				bufSize := 1024
				buffer := make([]byte, bufSize)
				length, _ := req.Body.Read(buffer)
				assert.Equal(t, req.Header.Get(tc.typeOrg), tc.idOrg)
				assert.Equal(t, req.Header.Get("Authorization"), tc.typeAuthorization+" "+tc.token)
				var request yandextrackerhttpclient.RequestComment
				json.Unmarshal(buffer[:length], &request)
				assert.Equal(t, tc.text, request.Text)
				rw.Write([]byte(`OK`))
			}))
			defer server.Close()
			client := server.Client()
			yandexTrackerHttpClient := yandextrackerhttpclient.NewYandexTracketHttpClient(server.URL,
				server.URL, tc.idOrg, tc.typeOrg,
				tc.typeAuthorization, tc.token, client)
			commentCreator := NewYandexTrackerCommentCreator(yandexTrackerHttpClient)
			commentCreator.CreateComment(context.Background(), tc.text, tc.taskId)
		})
	}
}
