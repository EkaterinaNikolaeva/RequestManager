package jirataskcreator

import (
	"context"

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

func (t JiraTaskCreator) CreateTask(ctx context.Context, requestedTask task.TaskCreateRequest) (task.TaskCreated, error) {
	link, id, err := t.jiraHttpClient.CreateTask(ctx, mapJiraIssueFromTask(requestedTask))
	return task.TaskCreated{
		Link:        link,
		Name:        requestedTask.Name,
		Description: requestedTask.Description,
		Type:        requestedTask.Type,
		Project:     requestedTask.Project,
		Id:          id,
	}, err
}

func mapJiraIssueFromTask(requestedTask task.TaskCreateRequest) jirahttpclient.JiraTaskCreationRequest {
	issue := jirahttpclient.JiraTaskCreationRequest{
		Fields: jirahttpclient.JiraTaskCreationFields{
			Project: jirahttpclient.JiraTaskCreationProject{
				Key: requestedTask.Project,
			},
			Summary:     requestedTask.Name,
			Description: requestedTask.Description,
			IssueType: jirahttpclient.JiraTaskCreationIssueType{
				Name: requestedTask.Type,
			},
		},
	}
	return issue
}
