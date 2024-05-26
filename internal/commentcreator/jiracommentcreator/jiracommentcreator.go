package jiracommentcreator

import (
	"context"
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

func (t JiraCommentCreator) CreateComment(ctx context.Context, text string, idTask string) error {
	err := t.jiraHttpClient.AddComment(ctx, text, idTask)
	if err != nil {
		return nil
	}
	log.Printf("Add comment in Jira: %s to task %s", text, idTask)
	return nil
}
