package service

import (
	"bytes"
	"context"
	"html/template"
	"log"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/task"
	errornotfound "github.com/EkaterinaNikolaeva/RequestManager/internal/storage/errors"
)

type StorageMsgTasks interface {
	GetIdTaskByMessage(ctx context.Context, msgId string) (string, error)
	GetIdMessageByTask(ctx context.Context, taskId string) (string, error)
	AddElement(ctx context.Context, msgId string, taskId string) error
	Finish()
}

type MessagesProvider interface {
	GetMessagesChannel() <-chan message.Message
}

type MessagesSender interface {
	SendMessage(message message.Message) error
}

type TaskCreator interface {
	CreateTask(task task.TaskCreateRequest) (task.TaskCreated, error)
}

type CommentCreator interface {
	CreateComment(text string, idTask string) error
}

type MessagesMatcher interface {
	MatchMessage(message message.Message) bool
}

type TaskFromMessagesCreator struct {
	messagesProvider    MessagesProvider
	messagesSender      MessagesSender
	taskCreator         TaskCreator
	messagesMatcher     MessagesMatcher
	messageReply        *template.Template
	taskDefaultProject  string
	taskDefaultType     string
	storageTaskMessages StorageMsgTasks
	commentCreator      CommentCreator
}

func NewTaskFromMessagesCreator(provider MessagesProvider, sender MessagesSender, matcher MessagesMatcher,
	taskCreator TaskCreator, messageDefaultReply *template.Template,
	taskDefaultProject string, taskDefaultType string, storage StorageMsgTasks, commentCreator CommentCreator) TaskFromMessagesCreator {
	return TaskFromMessagesCreator{
		messagesProvider:    provider,
		messagesSender:      sender,
		messagesMatcher:     matcher,
		messageReply:        messageDefaultReply,
		taskCreator:         taskCreator,
		taskDefaultProject:  taskDefaultProject,
		taskDefaultType:     taskDefaultType,
		storageTaskMessages: storage,
		commentCreator:      commentCreator,
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
			var isTask bool
			isTask = false
			if s.storageTaskMessages != nil {
				taskId, err := s.storageTaskMessages.GetIdTaskByMessage(ctx, msg.RootMessageId)
				if err != nil {
					_, ok := err.(errornotfound.NotFoundError)
					if !ok {
						log.Printf("error when get task id by msg %s: %q", msg.RootMessageId, err)
					}
				} else if !msg.Author.IsBot {
					isTask = true
					log.Printf("get message in thread %s by task %s", msg.RootMessageId, taskId)
					err := s.commentCreator.CreateComment("New msg in thread: "+msg.MessageText, taskId)
					if err != nil {
						log.Printf("error when add comment in thread %q", err)
					}
				}
			}
			if !msg.Author.IsBot && s.messagesMatcher.MatchMessage(msg) && !isTask {
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
				var reply bytes.Buffer
				err = s.messageReply.Execute(&reply, task)
				if err != nil {
					log.Printf("error when execute reply template %q", err)
					continue
				}
				err = s.messagesSender.SendMessage(
					message.Message{MessageText: reply.String(),
						ChannelId:     msg.ChannelId,
						RootMessageId: msg.RootMessageId})
				if err != nil {
					log.Printf("error when send reply %q", err)
				}
				if s.storageTaskMessages != nil {
					err = s.storageTaskMessages.AddElement(ctx, msg.RootMessageId, task.Id)
					if err != nil {
						log.Printf("error when try add element to storage %q", err)
					}
				}
			}

		}
	}
}
