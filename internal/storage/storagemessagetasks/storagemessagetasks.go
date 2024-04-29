package storagemessagetasks

import (
	"context"
	"log"
	"sync"
)

type StorageMsgTasksStupid struct {
	TaskByMessage map[string]string
	MessageByTask map[string]string
	mutex         sync.Mutex
}

func NewStorageMsgTasksStupid() StorageMsgTasksStupid {
	taskByMessage := make(map[string]string)
	messageByTask := make(map[string]string)
	return StorageMsgTasksStupid{
		TaskByMessage: taskByMessage,
		MessageByTask: messageByTask,
	}
}

func (s *StorageMsgTasksStupid) GetIdTaskByMessage(ctx context.Context, msgId string) (string, bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	value, ok := s.TaskByMessage[msgId]
	return value, ok, nil
}

func (s *StorageMsgTasksStupid) GetIdMessageByTask(ctx context.Context, taskId string) (string, bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	value, ok := s.MessageByTask[taskId]
	return value, ok, nil
}

func (s *StorageMsgTasksStupid) AddElement(ctx context.Context, msgId string, taskId string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	log.Printf("Add task %s, msg %s to storage", taskId, msgId)
	s.MessageByTask[taskId] = msgId
	s.TaskByMessage[msgId] = taskId
	return nil
}

func (s *StorageMsgTasksStupid) Finish() {

}
