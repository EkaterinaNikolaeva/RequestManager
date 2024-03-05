package jiratasks

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/task"
)

func basicAuth(username string, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

type JiraHttpClient struct {
	url               string
	httpClient        *http.Client
	authorizationCode string
}

func NewJiraHttpClient(httpClient *http.Client, url string, username string, password string) JiraHttpClient {
	return JiraHttpClient{
		url:               url,
		httpClient:        httpClient,
		authorizationCode: basicAuth(username, password),
	}
}

type JiraTaskCreationRequest struct {
	Fields JiraTaskCreationFields `json:"fields"`
}

type JiraTaskCreationFields struct {
	Project     JiraTaskCreationProject   `json:"project"`
	Summary     string                    `json:"summary"`
	Description string                    `json:"description"`
	IssueType   JiraTaskCreationIssueType `json:"issuetype"`
}

type JiraTaskCreationProject struct {
	Key string `json:"key"`
}

type JiraTaskCreationIssueType struct {
	Name string `json:"name"`
}

func (client *JiraHttpClient) CreateIssue(task task.Task) error {
	issue := JiraTaskCreationRequest{
		Fields: JiraTaskCreationFields{
			Project: JiraTaskCreationProject{
				Key: task.Project,
			},
			Summary:     task.Name,
			Description: task.Description,
			IssueType: JiraTaskCreationIssueType{
				Name: task.Type,
			},
		},
	}
	return client.makeRequestCreationTask(issue)
}

func (client *JiraHttpClient) makeRequestCreationTask(task JiraTaskCreationRequest) error {
	bytesRepresentation, err := json.Marshal(task)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", client.url+"/rest/api/2/issue/", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return err
	}
	log.Println(bytes.NewBuffer(bytesRepresentation))
	basic := "Basic " + client.authorizationCode
	req.Header.Add("Authorization", basic)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("Jira create task: %s", bytesResp)
	return nil
}
