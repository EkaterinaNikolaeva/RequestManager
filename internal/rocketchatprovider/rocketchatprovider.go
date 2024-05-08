package rocketchatprovider

import (
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

func (c *RocketChatPoviderClient) Login() error {
	request := ddpLogin{
		Token: c.token,
	}
	_, err := c.ddp.Call("login", request)
	if err != nil {
		return err
	}
	return nil
}

type messageExtractor struct {
	messageChannel chan message.Message
	operation      string
}

func (u messageExtractor) CollectionUpdate(collection, operation, id string, doc ddp.Update) {
	if operation == u.operation {
		document, _ := gabs.Consume(doc["args"])
		if document.Path("replies").Data() != nil {
			return
		}
		args, err := document.Children()
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

func (c *RocketChatPoviderClient) ProccessChannels(msgChannel chan message.Message) error {
	rawResponse, err := c.ddp.Call("rooms/get", map[string]int{
		"$date": 0,
	})
	if err != nil {
		return nil
	}
	document, _ := gabs.Consume(rawResponse.(map[string]interface{})["update"])
	log.Printf("%s", document)
	children, err := document.Children()
	if err != nil {
		return nil
	}
	for i := range children {
		id := children[i].Path("_id")
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

func (client *RocketChatPoviderClient) Run() error {
	err := client.Login()
	if err != nil {
		log.Fatalf("fatalf2 %q", err)
	}
	msgChannel := make(chan message.Message, 1024)
	err = client.ProccessChannels(msgChannel)
	if err != nil {
		log.Fatalf("fatalf3 %q", err)
	}
	for {
		select {
		case msg := <-msgChannel:
			log.Printf("THERE %s %s %s", msg.RootMessageId, msg.MessageText, msg.ChannelId)
		}
	}
}
