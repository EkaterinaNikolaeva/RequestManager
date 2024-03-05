package jirataskcreator

import (
	"github.com/EkaterinaNikolaeva/RequestManager/internal/jiratasks"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/task"
)

type JiraTaskCreator struct {
	jiraHttpClient jiratasks.JiraHttpClient
}

func NewJiraTaskCreator(jiraHttpClient jiratasks.JiraHttpClient) JiraTaskCreator {
	return JiraTaskCreator{
		jiraHttpClient: jiraHttpClient,
	}
}

func (t JiraTaskCreator) CreateTask(task task.Task) (task.Task, error) {
	return t.jiraHttpClient.CreateIssue(task)
}
