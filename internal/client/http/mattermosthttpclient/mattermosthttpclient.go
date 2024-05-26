package mattermosthttpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

func (client *MattermostHttpClient) CreatePost(ctx context.Context, post RequestPost) error {
	bytesRepresentation, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp marshal message for creation mattermost post")
	}
	req, err := http.NewRequestWithContext(ctx, "POST", client.mattermostBaseUrl+"/api/v4/posts", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp make new request message for creation mattermost post")
	}
	var bearer = "Bearer " + client.mattermostToken
	req.Header.Add("Authorization", bearer)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp do request for creation mattermost post")
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("http response: status code %d and status %s when attemp make request", resp.StatusCode, resp.Status)
	}
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp read bytes for creation mattermost post")
	}
	log.Printf("MattermostBot create post: %s", bytesResp)
	return nil
}
