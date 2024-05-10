package yandextrackertaskcreator

import (
	yandextrackerhttpclient "github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/yandextrackerclient"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/task"
)

type YandexTrackerTaskCreator struct {
	yandexTrackerHttpClient *yandextrackerhttpclient.YandexTracketHttpClient
}

func NewYandexTrackerTaskCreator(yandexTrackerHttpClient *yandextrackerhttpclient.YandexTracketHttpClient) YandexTrackerTaskCreator {
	return YandexTrackerTaskCreator{
		yandexTrackerHttpClient: yandexTrackerHttpClient,
	}
}

func (y YandexTrackerTaskCreator) CreateTask(task task.TaskCreateRequest) (task.TaskCreated, error) {
}
