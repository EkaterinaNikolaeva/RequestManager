package yandextrackerhttpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type YandexTracketHttpClient struct {
	host             string
	baseUrl          string
	idOrganization   string
	typeOrganization string
	tokenType        string
	token            string
	httpClient       *http.Client
}

func NewYandexTracketHttpClient(host string, baseUrl string, idOrganization string, typeOrganization string, tokenType string, token string, client *http.Client) YandexTracketHttpClient {
	return YandexTracketHttpClient{
		host:             host,
		baseUrl:          baseUrl,
		idOrganization:   idOrganization,
		typeOrganization: typeOrganization,
		tokenType:        tokenType,
		token:            token,
		httpClient:       client,
	}
}

func (client *YandexTracketHttpClient) getIssueLink(response ResponseTask) string {
	link := client.baseUrl + "/" + response.Key
	return link
}

func (client *YandexTracketHttpClient) addHeaders(req *http.Request) {
	req.Header.Add("Authorization", client.tokenType+" "+client.token)
	req.Header.Add("Host", client.host)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(client.typeOrganization, client.idOrganization)

}

func (client *YandexTracketHttpClient) CreateTask(ctx context.Context, task RequestTask) (string, string, error) {
	bytesRepresentation, err := json.Marshal(task)
	log.Printf("%s", bytesRepresentation)
	if err != nil {
		return "", "", fmt.Errorf(err.Error() + " when attemp create yandex track marshal task")
	}
	req, err := http.NewRequestWithContext(ctx, "POST", client.host+"/v2/issues/", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return "", "", fmt.Errorf(err.Error() + " when attemp new request for create yandex tracker task")
	}
	client.addHeaders(req)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf(err.Error() + " when attemp do http request for create yandex tracker task")
	}
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf(err.Error() + " when attemp create yandex tracker task")
	}
	log.Printf("Yandex Tracker create task: %s", bytesResp)
	var response ResponseTask
	err = json.Unmarshal(bytesResp, &response)
	if err != nil {
		return "", "", fmt.Errorf(err.Error() + " when attemp create yandex tracker task")
	}
	link := client.getIssueLink(response)
	return link, response.Key, nil
}

func (client *YandexTracketHttpClient) AddComment(ctx context.Context, text string, idIssue string) error {
	comment := RequestComment{
		Text: text,
	}
	bytesRepresentation, err := json.Marshal(comment)
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp marshal yandex tracker comment")
	}
	req, err := http.NewRequestWithContext(ctx, "POST", client.host+"/v2/issues/"+idIssue+"/comments", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp new request for create yandex tracker comment")
	}
	client.addHeaders(req)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp do http request for create yandex tracker comment")
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp readALl response body creation comment yandex tracker")
	}
	return nil
}
