package mattermosthttpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/mattermostmessages"
)

type MattermostHttpClient struct {
	httpClient        *http.Client
	mattermostToken   string
	mattermostBaseUrl string
}

func NewHttpClient(client *http.Client, token string, baseUrl string) *MattermostHttpClient {
	return &MattermostHttpClient{
		httpClient:        client,
		mattermostToken:   token,
		mattermostBaseUrl: baseUrl,
	}
}

func (client *MattermostHttpClient) CreatePost(post mattermostmessages.RequestPost) error {
	bytesRepresentation, err := json.Marshal(post)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", client.mattermostBaseUrl+"/api/v4/posts", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return err
	}
	var bearer = "Bearer " + client.mattermostToken
	req.Header.Add("Authorization", bearer)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("MattermostBot create post: %s", bytesResp)
	return nil
}
