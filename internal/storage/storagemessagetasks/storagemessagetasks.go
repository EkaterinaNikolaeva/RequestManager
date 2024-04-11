package storagemessagetasks

import "log"

type StorageMsgTasksStupid struct {
	TaskByMessage map[string]string
	MessageByTask map[string]string
}

func NewStorageMsgTasksStupid() StorageMsgTasksStupid {
	taskByMessage := make(map[string]string)
	messageByTask := make(map[string]string)
	return StorageMsgTasksStupid{taskByMessage, messageByTask}
}

func (s *StorageMsgTasksStupid) GetIdTaskByMessage(msgId string) (string, bool) {
	value, ok := s.TaskByMessage[msgId]
	return value, ok
}

func (s *StorageMsgTasksStupid) GetIdMessageByTask(taskId string) (string, bool) {
	value, ok := s.MessageByTask[taskId]
	return value, ok
}

func (s *StorageMsgTasksStupid) AddElement(msgId string, taskId string) {
	log.Printf("Add task %s, msg %s to storage", taskId, msgId)
	s.MessageByTask[taskId] = msgId
	s.TaskByMessage[msgId] = taskId
}
