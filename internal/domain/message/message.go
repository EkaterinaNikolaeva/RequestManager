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

func NewMessage(text string, channelId string, rootId string, author MessageAuthor) Message {
	return Message{
		MessageText:   text,
		ChannelId:     channelId,
		RootMessageId: rootId,
		Author:        author,
	}
}

func NewMessageAuthor(id string, isBot ...bool) MessageAuthor {
	if len(isBot) > 0 {
		return MessageAuthor{
			Id:    id,
			IsBot: isBot[0],
		}
	}
	return MessageAuthor{
		Id:    id,
		IsBot: false,
	}
}
