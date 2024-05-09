package rocketchatprovider

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/domain/message"
	"github.com/Jeffail/gabs"
	"github.com/gopackage/ddp"
)

type RocketChatPoviderClient struct {
	ddp   *ddp.Client
	token string
}

func NewRocketChatPoviderClient(host string, token string) (*RocketChatPoviderClient, error) {
	wsUrl := fmt.Sprintf("ws://%v/websocket", host)
	client := RocketChatPoviderClient{
		ddp:   ddp.NewClient(wsUrl, (&url.URL{Host: host}).String()),
		token: token,
	}
	err := client.ddp.Connect()
	if err != nil {
		return nil, err
	}
	return &client, err
}

type ddpLogin struct {
	Token string `json:"resume"`
}

func (c *RocketChatPoviderClient) login() error {
	_, err := c.ddp.Call("login", ddpLogin{
		Token: c.token,
	})
	return err
}

type messageExtractor struct {
	messageChannel chan message.Message
	operation      string
}

func (u messageExtractor) CollectionUpdate(collection, operation, id string, doc ddp.Update) {
	if operation == u.operation {
		allArgs, _ := gabs.Consume(doc["args"])
		if allArgs.Path("replies").Data() != nil {
			return
		}
		args, err := allArgs.Children()
		if err != nil {
			return
		}
		for _, arg := range args {
			id, _ := arg.Path("_id").Data().(string)
			rid, _ := arg.Path("rid").Data().(string)
			msg, _ := arg.Path("msg").Data().(string)
			tmid, _ := arg.Path("tmid").Data().(string)
			username, _ := arg.Path("u.username").Data().(string)
			if tmid != "" {
				id = tmid
			}
			u.messageChannel <- message.Message{
				MessageText:   msg,
				ChannelId:     rid,
				RootMessageId: id,
				Author:        message.MessageAuthor{Id: username},
			}

		}
	}
}

func (c *RocketChatPoviderClient) proccessChannels(msgChannel chan message.Message) error {
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
		log.Println(id)
		err = c.ddp.Sub("stream-room-messages", id.Data(), true)
		if err != nil {
			return nil
		}
		if i == 0 {
			c.ddp.CollectionByName("stream-room-messages").AddUpdateListener(messageExtractor{msgChannel, "update"})
		}
	}
	return nil
}

func (client *RocketChatPoviderClient) Run(ctx context.Context) error {
	err := client.login()
	if err != nil {
		return fmt.Errorf("error when rocketchat login: %q", err)
	}
	msgChannel := make(chan message.Message, 1024)
	err = client.proccessChannels(msgChannel)
	if err != nil {
		return fmt.Errorf("error when rocketchat proccess msgs: %q", err)
	}
	for {
		select {
		case <-ctx.Done():
			log.Printf("ctx is done, stop rocketchat client")
			client.ddp.Close()
			return nil
		}
	}
}
