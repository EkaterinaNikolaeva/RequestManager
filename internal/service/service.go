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
}

type MessagesSender interface {
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
	messagesSender   MessagesSender
	taskCreator      TaskCreator
	messagesMatcher  MessagesMatcher
	messageReply     string
}

func NewTaskFromMessagesCreator(provider MessagesProvider, sender MessagesSender, matcher MessagesMatcher, messageStandardReply string) TaskFromMessagesCreator {
	return TaskFromMessagesCreator{
		messagesProvider: provider,
		messagesSender:   sender,
		messagesMatcher:  matcher,
		messageReply:     messageStandardReply,
	}
}

func (s TaskFromMessagesCreator) Run() {
	messagesChannel := s.messagesProvider.GetMessagesChannel()
	log.Println("Server started!")
	for {
		select {
		case msg := <-messagesChannel:
			if !msg.Author.IsBot && s.messagesMatcher.MatchMessage(msg) {
				s.messagesSender.SendMessage(
					message.Message{MessageText: s.messageReply,
						ChannelId:     msg.ChannelId,
						RootMessageId: msg.RootMessageId})
			}
		}
	}
}
