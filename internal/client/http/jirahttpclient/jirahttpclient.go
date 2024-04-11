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

func (client *JiraHttpClient) getIssueLink(response jiratasks.JiraTaskCreationResponse, task jiratasks.JiraTaskCreationRequest) string {
	link := client.baseUrl + "/projects/" + task.Fields.Project.Key + "/issues/" + response.Key
	return link
}

func (client *JiraHttpClient) CreateTask(task jiratasks.JiraTaskCreationRequest) (string, string, error) {
	bytesRepresentation, err := json.Marshal(task)
	if err != nil {
		return "", "", fmt.Errorf(err.Error() + " when attemp create jira issue marshal task")
	}
	req, err := http.NewRequest("POST", client.baseUrl+"/rest/api/2/issue/", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return "", "", fmt.Errorf(err.Error() + " when attemp new request for create jira task")
	}
	req.Header.Add("Authorization", "Basic "+client.authorizationCode)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf(err.Error() + " when attemp do http request for create jira task")
	}
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf(err.Error() + " when attemp create jira task")
	}
	log.Printf("Jira create task: %s", bytesResp)
	var response jiratasks.JiraTaskCreationResponse
	err = json.Unmarshal(bytesResp, &response)
	if err != nil {
		return "", "", fmt.Errorf(err.Error() + " when attemp create jira post")
	}
	link := client.getIssueLink(response, task)
	return link, response.Key, nil
}

func (client *JiraHttpClient) AddComment(text string, idIssue string) error {
	comment := jiratasks.JiraCommentRequest{
		Body: text,
	}
	bytesRepresentation, err := json.Marshal(comment)
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp marshal jira comment")
	}
	req, err := http.NewRequest("POST", client.baseUrl+"/rest/api/2/issue/"+idIssue+"/comment", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp new request for create jira comment")
	}
	req.Header.Add("Authorization", "Basic "+client.authorizationCode)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp do http request for create jira post")
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(err.Error() + " when attemp readALl response body creation comment")
	}
	return nil
}
