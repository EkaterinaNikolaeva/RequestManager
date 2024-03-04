package service

import (
	"log"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/message"
)

type Task struct {
	Text string
}

type MessagesProvider interface {
	GetMessagesChannel() <-chan message.Message
	SendMessage(message message.Message) error
}

type TaskCreator interface {
	CreateTask(task Task)
}

type MessagesMatcher interface {
	MatchMessage(message message.Message) bool
}

type TaskFromMessagesCreator struct {
	messagesProvider MessagesProvider
	taskCreator      TaskCreator
	messagesMatcher  MessagesMatcher
}

func NewTaskFromMessagesCreator(provider MessagesProvider, matcher MessagesMatcher) TaskFromMessagesCreator {
	return TaskFromMessagesCreator{
		messagesProvider: provider,
		messagesMatcher:  matcher,
	}
}
func (s TaskFromMessagesCreator) Run() {
	messagesChannel := s.messagesProvider.GetMessagesChannel()
	log.Println("Server started!")
	for {
		msg := <-messagesChannel
		if s.messagesMatcher.MatchMessage(msg) {
			s.messagesProvider.SendMessage(
				message.Message{MessageText: "Make issue! Link: ",
					ChannelId:     msg.ChannelId,
					RootMessageId: msg.RootMessageId})
		}
	}
}
