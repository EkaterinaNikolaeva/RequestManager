package storagemessagetasks

import (
	"context"
	"log"
)

type StorageMsgTasksStupid struct {
	TaskByMessage map[string]string
	MessageByTask map[string]string
}

func NewStorageMsgTasksStupid() StorageMsgTasksStupid {
	taskByMessage := make(map[string]string)
	messageByTask := make(map[string]string)
	return StorageMsgTasksStupid{taskByMessage, messageByTask}
}

func (s *StorageMsgTasksStupid) GetIdTaskByMessage(msgId string, ctx context.Context) (string, bool, error) {
	value, ok := s.TaskByMessage[msgId]
	return value, ok, nil
}

func (s *StorageMsgTasksStupid) GetIdMessageByTask(taskId string, ctx context.Context) (string, bool, error) {
	value, ok := s.MessageByTask[taskId]
	return value, ok, nil
}

func (s *StorageMsgTasksStupid) AddElement(msgId string, taskId string, ctx context.Context) error {
	log.Printf("Add task %s, msg %s to storage", taskId, msgId)
	s.MessageByTask[taskId] = msgId
	s.TaskByMessage[msgId] = taskId
	return nil
}

func (s *StorageMsgTasksStupid) Finish() {

}
