package service

type Message struct {
	Message   string
	Chat      string
	MessageId string
}

type Task struct {
	Text string
}

type MessagesProvider interface {
	GetMessagesChannel() <-chan Message
	SendMessage(message Message) error
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
		select {
		case message := <-messagesChannel:
			s.messagesProvider.SendMessage(message)
			// if s.messagesMatcher.MatchMessage(message) {
			// 	// s.taskCreator.CreateTask()
			// }
		}

	}

}
