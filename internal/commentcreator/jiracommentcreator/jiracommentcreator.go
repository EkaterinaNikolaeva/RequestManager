package jiracommentcreator

import (
	"log"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/jirahttpclient"
)

type JiraCommentCreator struct {
	jiraHttpClient jirahttpclient.JiraHttpClient
}

func NewJiraCommentCreator(jiraHttpClient jirahttpclient.JiraHttpClient) JiraCommentCreator {
	return JiraCommentCreator{
		jiraHttpClient: jiraHttpClient,
	}
}

func (t JiraCommentCreator) CreateComment(text string, idTask string) error {
	err := t.jiraHttpClient.AddComment(text, idTask)
	if err != nil {
		return nil
	}
	log.Printf("Add comment in Jira: %s to task %s", text, idTask)
	return nil
}
