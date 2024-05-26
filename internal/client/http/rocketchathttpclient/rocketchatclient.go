package rocketchathttpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type RocketChatHttpClient struct {
	httpClient        *http.Client
	rocketChatBaseUrl string
	rocketChatId      string
	rocketChatToken   string
}

func NewHttpClient(client *http.Client, id string, token string, baseUrl string) *RocketChatHttpClient {
	return &RocketChatHttpClient{
		httpClient:        client,
		rocketChatId:      id,
		rocketChatToken:   token,
		rocketChatBaseUrl: baseUrl,
	}
}

func (client *RocketChatHttpClient) SendMessage(ctx context.Context, msg RequestMessage) error {
	bytesRepresentation, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp marshal message for creation rocketchat post")
	}
	req, err := http.NewRequestWithContext(ctx, "POST", client.rocketChatBaseUrl+"/api/v1/chat.sendMessage", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp make new request message for creation rocketchat post")
	}
	req.Header.Add("X-User-Id", client.rocketChatId)
	req.Header.Add("X-Auth-Token", client.rocketChatToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp do request for creation rocketchat post")
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("http response rocketchat: status code %d and status %s when attemp make request", resp.StatusCode, resp.Status)
	}
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp read bytes for creation rocketchat post")
	}
	log.Printf("RocketChat send message: %s", bytesResp)
	return nil
}
