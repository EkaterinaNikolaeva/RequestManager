package yandextrackerhttpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	apiyandextracker "github.com/EkaterinaNikolaeva/RequestManager/internal/api/yandextracker"
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

func (client *YandexTracketHttpClient) getIssueLink(response apiyandextracker.ResponseTask) string {
	link := client.baseUrl + "/" + response.Key
	return link
}

func (client *YandexTracketHttpClient) CreateTask(task apiyandextracker.RequestTask) (string, string, error) {
	bytesRepresentation, err := json.Marshal(task)
	log.Printf("%s", bytesRepresentation)
	if err != nil {
		return "", "", fmt.Errorf(err.Error() + " when attemp create yandex track marshal task")
	}
	req, err := http.NewRequest("POST", client.host+"/v2/issues/", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return "", "", fmt.Errorf(err.Error() + " when attemp new request for create yandex tracker task")
	}
	req.Header.Add("Authorization", client.tokenType+" "+client.token)
	req.Header.Add("Host", client.host)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(client.typeOrganization, client.idOrganization)
	log.Printf("%s", req.Header)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf(err.Error() + " when attemp do http request for create yandex tracker task")
	}
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf(err.Error() + " when attemp create yandex tracker task")
	}
	log.Printf("Yandex Tracker create task: %s", bytesResp)
	var response apiyandextracker.ResponseTask
	err = json.Unmarshal(bytesResp, &response)
	if err != nil {
		return "", "", fmt.Errorf(err.Error() + " when attemp create jira post")
	}
	link := client.getIssueLink(response)
	return link, response.Key, nil
}
