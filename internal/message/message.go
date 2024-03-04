package message

type Message struct {
	MessageText   string
	ChannelId     string
	RootMessageId string
	Author        MessageAuthor
}

type MessageAuthor struct {
	Id    string
	IsBot bool
}
