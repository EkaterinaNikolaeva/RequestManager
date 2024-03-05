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

type JiraTaskCreationResponse struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

func (client *JiraHttpClient) CreateIssue(task task.Task) (task.Task, error) {
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
	link, err := client.makeRequestCreationTask(issue)
	task.Link = link
	return task, err
}

func (client *JiraHttpClient) getIssueLink(bytes []byte, task JiraTaskCreationRequest) (string, error) {
	var response JiraTaskCreationResponse
	err := json.Unmarshal(bytes, &response)
	if err != nil {
		return "", err
	}
	link := client.url + "/projects/" + task.Fields.Project.Key + "/issues/" + response.Key
	return link, nil
}

func (client *JiraHttpClient) makeRequestCreationTask(task JiraTaskCreationRequest) (string, error) {
	bytesRepresentation, err := json.Marshal(task)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", client.url+"/rest/api/2/issue/", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Basic "+client.authorizationCode)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Printf("Jira create task: %s", bytesResp)
	return client.getIssueLink(bytesResp, task)
}
