package service

import "log"

type Message struct {
	Message string
}

type Task struct {
	Text string
}

type MessagesProvider interface {
	GetMessagesChannel() <-chan Message
}

type TaskCreator interface {
	CreateTask(task Task)
}

type MessagesMatcher interface {
	MatchMessage(message Message)
}

type TaskFromMessagesCreator struct {
	messagesProvider MessagesProvider
	taskCreator      TaskCreator
	messagesMatcher  MessagesMatcher
}

func NewTaskFromMessagesCreator(provider MessagesProvider) TaskFromMessagesCreator {
	return TaskFromMessagesCreator{
		messagesProvider: provider,
	}
}
func (s TaskFromMessagesCreator) Run() {
	messagesChannel := s.messagesProvider.GetMessagesChannel()
	for {
		message := <-messagesChannel
		log.Println(message.Message)
		// if s.messagesMatcher.MatchMessage(message) {
		// 	// s.taskCreator.CreateTask()
		// }

	}

}
