package storageinmemory

import (
	"context"
	"log"
	"sync"

	errornotfound "github.com/EkaterinaNikolaeva/RequestManager/internal/storage/errors"
)

type StorageMsgTasksInMemory struct {
	TaskByMessage map[string]string
	MessageByTask map[string]string
	rwMutex       sync.RWMutex
}

func NewStorageMsgTasksInMemory() StorageMsgTasksInMemory {
	taskByMessage := make(map[string]string)
	messageByTask := make(map[string]string)
	return StorageMsgTasksInMemory{
		TaskByMessage: taskByMessage,
		MessageByTask: messageByTask,
	}
}

func (s *StorageMsgTasksInMemory) GetIdTaskByMessage(ctx context.Context, msgId string) (string, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	value, ok := s.TaskByMessage[msgId]
	if !ok {
		return value, errornotfound.NewNotFoundError()
	}
	return value, nil
}

func (s *StorageMsgTasksInMemory) GetIdMessageByTask(ctx context.Context, taskId string) (string, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	value, ok := s.MessageByTask[taskId]
	if !ok {
		return value, errornotfound.NewNotFoundError()
	}
	return value, nil
}

func (s *StorageMsgTasksInMemory) AddElement(ctx context.Context, msgId string, taskId string) error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	log.Printf("Add task %s, msg %s to storage", taskId, msgId)
	s.MessageByTask[taskId] = msgId
	s.TaskByMessage[msgId] = taskId
	return nil
}

func (s *StorageMsgTasksInMemory) Finish() {

}
