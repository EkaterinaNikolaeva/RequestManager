package jirataskcreator

import (
	"log"

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

func (t JiraTaskCreator) CreateTask(task task.Task) {
	err := t.jiraHttpClient.CreateIssue(task)
	if err != nil {
		log.Printf("%q", err)
	}
}
