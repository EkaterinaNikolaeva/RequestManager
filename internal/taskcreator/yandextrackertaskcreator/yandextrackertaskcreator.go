package yandextrackertaskcreator

import (
	"github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/yandextrackerhttpclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/task"
)

type YandexTrackerTaskCreator struct {
	yandexTrackerHttpClient yandextrackerhttpclient.YandexTracketHttpClient
}

func NewYandexTrackerTaskCreator(yandexTrackerHttpClient yandextrackerhttpclient.YandexTracketHttpClient) YandexTrackerTaskCreator {
	return YandexTrackerTaskCreator{
		yandexTrackerHttpClient: yandexTrackerHttpClient,
	}
}

func (y YandexTrackerTaskCreator) CreateTask(requestedTask task.TaskCreateRequest) (task.TaskCreated, error) {
	link, id, err := y.yandexTrackerHttpClient.CreateTask(mapJiraIssueFromTask(requestedTask))
	return task.TaskCreated{
		Link:        link,
		Name:        requestedTask.Name,
		Description: requestedTask.Description,
		Type:        requestedTask.Type,
		Project:     requestedTask.Project,
		Id:          id,
	}, err
}

func mapJiraIssueFromTask(requestedTask task.TaskCreateRequest) yandextrackerhttpclient.RequestTask {
	return yandextrackerhttpclient.RequestTask{
		Summary:     requestedTask.Name,
		Description: requestedTask.Description,
		Queue:       requestedTask.Project,
		Type:        requestedTask.Type,
	}
}
