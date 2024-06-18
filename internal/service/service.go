package service

import (
	"bytes"
	"context"
	"html/template"
	"log"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/config"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/task"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/storage/storageerrors"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=service .

type StorageMsgTasks interface {
	GetIdTaskByMessage(ctx context.Context, msgId string) (string, error)
	GetIdMessageByTask(ctx context.Context, taskId string) (string, error)
	AddElement(ctx context.Context, msgId string, taskId string) error
	Finish()
}

type MessagesProvider interface {
	GetMessagesChannel() <-chan message.Message
	Run(ctx context.Context)
}

type MessagesSender interface {
	SendMessage(ctx context.Context, msg message.Message) error
}

type TaskCreator interface {
	CreateTask(ctx context.Context, taskRequest task.TaskCreateRequest) (task.TaskCreated, error)
}

type CommentCreator interface {
	CreateComment(ctx context.Context, text string, idTask string) error
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
	defaultTaskName     *template.Template
	Messenger           config.Messenger
	taskDefaultProject  string
	taskDefaultType     string
	enableMsgThreating  bool
	storageTaskMessages StorageMsgTasks
	commentCreator      CommentCreator
}

func NewTaskFromMessagesCreator(provider MessagesProvider, sender MessagesSender, matcher MessagesMatcher,
	taskCreator TaskCreator, messageDefaultReply *template.Template,
	taskDefaultProject string, taskDefaultType string, enableMsgThreating bool, storage StorageMsgTasks,
	commentCreator CommentCreator, messenger config.Messenger, defaultTaskName *template.Template) TaskFromMessagesCreator {
	return TaskFromMessagesCreator{
		messagesProvider:    provider,
		messagesSender:      sender,
		messagesMatcher:     matcher,
		messageReply:        messageDefaultReply,
		taskCreator:         taskCreator,
		taskDefaultProject:  taskDefaultProject,
		taskDefaultType:     taskDefaultType,
		enableMsgThreating:  enableMsgThreating,
		storageTaskMessages: storage,
		commentCreator:      commentCreator,
		Messenger:           messenger,
		defaultTaskName:     defaultTaskName,
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
			if s.enableMsgThreating {
				taskId, err := s.storageTaskMessages.GetIdTaskByMessage(ctx, msg.RootMessageId)
				if err != nil {
					_, ok := err.(storageerrors.NotFoundError)
					if !ok {
						log.Printf("error when get task id by msg %s: %q", msg.RootMessageId, err)
					}
				} else if !msg.Author.IsBot {
					isTask = true
					log.Printf("get message in thread %s by task %s", msg.RootMessageId, taskId)
					err := s.commentCreator.CreateComment(ctx, "New msg in thread: "+msg.MessageText, taskId)
					if err != nil {
						log.Printf("error when add comment in thread %q", err)
					}
				}
			}
			if !msg.Author.IsBot && s.messagesMatcher.MatchMessage(msg) && !isTask {
				var taskName bytes.Buffer
				err := s.defaultTaskName.Execute(&taskName, s)
				if err != nil {
					log.Printf("error when execute task name template %q", err)
					continue
				}
				task, err := s.taskCreator.CreateTask(ctx,
					task.TaskCreateRequest{
						Name:        taskName.String(),
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
				err = s.messagesSender.SendMessage(ctx,
					message.Message{MessageText: reply.String(),
						ChannelId:     msg.ChannelId,
						RootMessageId: msg.RootMessageId})
				if err != nil {
					log.Printf("error when send reply %q", err)
				}
				if s.enableMsgThreating {
					err = s.storageTaskMessages.AddElement(ctx, msg.RootMessageId, task.Id)
					if err != nil {
						log.Printf("error when try add element to storage %q", err)
					}
				}
			}

		}
	}
}
