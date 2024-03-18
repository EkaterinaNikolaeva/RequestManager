package service

import (
	"context"
	"log"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/task"
)

type MessagesProvider interface {
	GetMessagesChannel() <-chan message.Message
}

type MessagesSender interface {
	SendMessage(message message.Message) error
}

type TaskCreator interface {
	CreateTask(task task.TaskCreateRequest) (task.TaskCreated, error)
}

type MessagesMatcher interface {
	MatchMessage(message message.Message) bool
}

type TaskFromMessagesCreator struct {
	messagesProvider   MessagesProvider
	messagesSender     MessagesSender
	taskCreator        TaskCreator
	messagesMatcher    MessagesMatcher
	messageReply       string
	taskDefaultProject string
	taskDefaultType    string
}

func NewTaskFromMessagesCreator(provider MessagesProvider, sender MessagesSender, matcher MessagesMatcher,
	taskCreator TaskCreator, messageDefaultReply string,
	taskDefaultProject string, taskDefaultType string) TaskFromMessagesCreator {
	return TaskFromMessagesCreator{
		messagesProvider:   provider,
		messagesSender:     sender,
		messagesMatcher:    matcher,
		messageReply:       messageDefaultReply,
		taskCreator:        taskCreator,
		taskDefaultProject: taskDefaultProject,
		taskDefaultType:    taskDefaultType,
	}
}

func (s TaskFromMessagesCreator) Run(ctx context.Context) {
	messagesChannel := s.messagesProvider.GetMessagesChannel()
	log.Println("Server started!")
	for {
		select {
		case <-ctx.Done():
			log.Printf("ctx is done, stop service task from message creation")
			return
		case msg := <-messagesChannel:
			if !msg.Author.IsBot && s.messagesMatcher.MatchMessage(msg) {
				task, err := s.taskCreator.CreateTask(
					task.TaskCreateRequest{
						Name:        "From mattermost",
						Description: msg.MessageText,
						Project:     s.taskDefaultProject,
						Type:        s.taskDefaultType,
					})
				if err != nil {
					log.Printf("error when create task %q", err)
					continue
				}
				err = s.messagesSender.SendMessage(
					message.Message{MessageText: s.messageReply + task.Link,
						ChannelId:     msg.ChannelId,
						RootMessageId: msg.RootMessageId})
				if err != nil {
					log.Printf("error when send reply %q", err)
				}
			}
		}
	}
}
