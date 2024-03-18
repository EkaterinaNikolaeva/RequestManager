package jirahttpclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/jiratasks"
)

func basicAuth(username string, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

type JiraHttpClient struct {
	baseUrl           string
	httpClient        *http.Client
	authorizationCode string
}

func NewJiraHttpClient(httpClient *http.Client, url string, username string, password string) JiraHttpClient {
	return JiraHttpClient{
		baseUrl:           url,
		httpClient:        httpClient,
		authorizationCode: basicAuth(username, password),
	}
}

func (client *JiraHttpClient) getIssueLink(response jiratasks.JiraTaskCreationResponse, task jiratasks.JiraTaskCreationRequest) (string, error) {
	link := client.baseUrl + "/projects/" + task.Fields.Project.Key + "/issues/" + response.Key
	return link, nil
}

func (client *JiraHttpClient) CreateTask(task jiratasks.JiraTaskCreationRequest) (string, error) {
	bytesRepresentation, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf(err.Error() + " when attemp create jira issue marshal task")
	}
	req, err := http.NewRequest("POST", client.baseUrl+"/rest/api/2/issue/", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return "", fmt.Errorf(err.Error() + " when attemp new request for create mattermost post")
	}
	req.Header.Add("Authorization", "Basic "+client.authorizationCode)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf(err.Error() + " when attemp do http request for create mattermost post")
	}
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf(err.Error() + " when attemp create mattermost post")
	}
	log.Printf("Jira create task: %s", bytesResp)
	var response jiratasks.JiraTaskCreationResponse
	err = json.Unmarshal(bytesResp, &response)
	if err != nil {
		return "", fmt.Errorf(err.Error() + " when attemp create mattermost post")
	}
	return client.getIssueLink(response, task)
}
