package rocketchatprovider

import (
	"context"
	"log"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
	"github.com/Jeffail/gabs"
	"github.com/gopackage/ddp"
)

type RocketChatPovider struct {
	ddp        *ddp.Client
	token      string
	id         string
	msgChannel chan message.Message
}

func NewRocketChatPovider(ddpClient *ddp.Client, id string, token string) (*RocketChatPovider, error) {
	msgChannel := make(chan message.Message, 1024)
	provider := RocketChatPovider{
		ddp:        ddpClient,
		token:      token,
		id:         id,
		msgChannel: msgChannel,
	}
	err := provider.ddp.Connect()
	if err != nil {
		return nil, err
	}
	return &provider, err
}

type ddpLogin struct {
	Token string `json:"resume"`
}

func (c *RocketChatPovider) login() error {
	_, err := c.ddp.Call("login", ddpLogin{
		Token: c.token,
	})
	return err
}

type messageExtractor struct {
	messageChannel chan message.Message
	operation      string
	botId          string
}

func (m messageExtractor) CollectionUpdate(collection, operation, id string, doc ddp.Update) {
	if operation == m.operation {
		allArgs, _ := gabs.Consume(doc["args"])
		if allArgs.Path("replies").Data() != nil {
			return
		}
		args, err := allArgs.Children()
		if err != nil {
			return
		}
		for _, arg := range args {
			var rootMsgId string
			id, ok := arg.Path("_id").Data().(string)
			if !ok {
				continue
			}
			rootMsgId = id
			rid, ok := arg.Path("rid").Data().(string)
			if !ok {
				continue
			}
			msg, ok := arg.Path("msg").Data().(string)
			if !ok {
				continue
			}
			userId, ok := arg.Path("u._id").Data().(string)
			if !ok {
				continue
			}
			tmid, ok := arg.Path("tmid").Data().(string)
			if ok && tmid != "" {
				rootMsgId = tmid
			}
			m.messageChannel <- message.Message{
				MessageText:   msg,
				ChannelId:     rid,
				RootMessageId: rootMsgId,
				Author:        message.MessageAuthor{Id: userId, IsBot: userId == m.botId},
			}

		}
	}
}

func (c *RocketChatPovider) proccessChannels(msgChannel chan message.Message) error {
	rawResponse, err := c.ddp.Call("rooms/get", map[string]int{
		"$date": 0,
	})
	if err != nil {
		return nil
	}
	allChannels, _ := gabs.Consume(rawResponse.(map[string]interface{})["update"])
	channells, err := allChannels.Children()
	if err != nil {
		return nil
	}
	for i := range channells {
		id := channells[i].Path("_id")
		err = c.ddp.Sub("stream-room-messages", id.Data(), true)
		if err != nil {
			return nil
		}
		if i == 0 {
			c.ddp.CollectionByName("stream-room-messages").AddUpdateListener(messageExtractor{msgChannel, "update", c.id})
		}
	}
	return nil
}

func (p *RocketChatPovider) Run(ctx context.Context) {
	err := p.login()
	if err != nil {
		log.Fatalf("error when rocketchat login: %q", err)
	}
	err = p.proccessChannels(p.msgChannel)
	if err != nil {
		log.Fatalf("error when rocketchat proccess msgs: %q", err)
	}
	for {
		select {
		case <-ctx.Done():
			log.Printf("ctx is done, stop rocketchat provider")
			p.ddp.Close()
			return
		}
	}
}

func (c *RocketChatPovider) GetMessagesChannel() <-chan message.Message {
	return c.msgChannel
}
