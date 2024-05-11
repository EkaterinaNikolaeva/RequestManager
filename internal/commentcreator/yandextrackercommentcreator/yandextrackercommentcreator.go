package yandextrackercommentcreator

import (
	"log"

	yandextrackerhttpclient "github.com/EkaterinaNikolaeva/RequestManager/internal/client/http/yandextrackerclient"
)

type YandexTrackerCommentCreator struct {
	yandexTrackerHttpClient yandextrackerhttpclient.YandexTracketHttpClient
}

func NewYandexTrackerCommentCreator(yandexTrackerHttpClient yandextrackerhttpclient.YandexTracketHttpClient) YandexTrackerCommentCreator {
	return YandexTrackerCommentCreator{
		yandexTrackerHttpClient: yandexTrackerHttpClient,
	}
}

func (y YandexTrackerCommentCreator) CreateComment(text string, idTask string) error {
	err := y.yandexTrackerHttpClient.AddComment(text, idTask)
	if err != nil {
		return nil
	}
	log.Printf("Add comment in YandexTracker: %s to task %s", text, idTask)
	return nil
}
