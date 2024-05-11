package jirataskcreator

import (
	"github.com/EkaterinaNikolaeva/RequestManager/internal/api/jira/jiratasks"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/jirahttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/task"
)

type JiraTaskCreator struct {
	jiraHttpClient jirahttpclient.JiraHttpClient
}

func NewJiraTaskCreator(jiraHttpClient jirahttpclient.JiraHttpClient) JiraTaskCreator {
	return JiraTaskCreator{
		jiraHttpClient: jiraHttpClient,
	}
}

func (t JiraTaskCreator) CreateTask(requestedTask task.TaskCreateRequest) (task.TaskCreated, error) {
	link, id, err := t.jiraHttpClient.CreateTask(mapJiraIssueFromTask(requestedTask))
	return task.TaskCreated{
		Link:        link,
		Name:        requestedTask.Name,
		Description: requestedTask.Description,
		Type:        requestedTask.Type,
		Project:     requestedTask.Project,
		Id:          id,
	}, err
}

func mapJiraIssueFromTask(requestedTask task.TaskCreateRequest) jiratasks.JiraTaskCreationRequest {
	issue := jiratasks.JiraTaskCreationRequest{
		Fields: jiratasks.JiraTaskCreationFields{
			Project: jiratasks.JiraTaskCreationProject{
				Key: requestedTask.Project,
			},
			Summary:     requestedTask.Name,
			Description: requestedTask.Description,
			IssueType: jiratasks.JiraTaskCreationIssueType{
				Name: requestedTask.Type,
			},
		},
	}
	return issue
}
