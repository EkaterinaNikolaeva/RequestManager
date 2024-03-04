package mattermostmessages

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/message"
)

type HttpClient http.Client

func NewHttpClient(client *http.Client) *HttpClient {
	return (*HttpClient)(client)
}

type RequestPost struct {
	ChannelId string                 `json:"channel_id"`
	Message   string                 `json:"message"`
	RootId    string                 `json:"root_id,omitempty"`
	FileIds   []string               `json:"file_ids,omitempty"`
	Props     interface{}            `json:"props,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type ResponsePost struct {
	Id            string                 `json:"id,omitempty"`
	CreateAt      int                    `json:"create_at,omitempty"`
	UpdateAt      int                    `json:"update_at,omitempty"`
	DeleteAt      int                    `json:"delete_at,omitempty"`
	EditAt        int                    `json:"edit_at,omitempty"`
	UserId        string                 `json:"user_id,omitempty"`
	ChannelId     string                 `json:"channel_id,omitempty"`
	RootId        string                 `json:"root_id,omitempty"`
	OriginalId    string                 `json:"original_id,omitempty"`
	Message       string                 `json:"message,omitempty"`
	Type          string                 `json:"type,omitempty"`
	Props         map[string]interface{} `json:"props,omitempty"`
	Hashtag       string                 `json:"hashtag,omitempty"`
	PendingPostId string                 `json:"pending_post_id,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

func getFrom(post ResponsePost) string {
	rootId := post.Id
	if post.RootId != "" {
		rootId = post.RootId
	}
	return rootId
}

func checkMessageFromBot(post ResponsePost) bool {
	props := post.Props
	fromBot, ok := props["from_bot"]
	isBot, isBool := fromBot.(bool)
	return ok && isBool && isBot || fromBot == "true"
}
func GetMessage(bytes string) (message.Message, error) {
	var post ResponsePost
	err := json.Unmarshal([]byte(bytes), &post)
	if err != nil {
		return message.Message{}, err
	}
	return message.Message{
		MessageText:   post.Message,
		ChannelId:     post.ChannelId,
		RootMessageId: getFrom(post),
		Author:        message.MessageAuthor{Id: post.UserId, IsBot: checkMessageFromBot(post)},
	}, nil
}

func (client *HttpClient) SendMessage(message message.Message, bot bot.MattermostBot) error {
	post := RequestPost{
		Message:   message.MessageText,
		ChannelId: message.ChannelId,
		RootId:    message.RootMessageId,
	}
	return client.CreatePost(post, bot.MattermostHttp, bot)
}

func (client *HttpClient) CreatePost(post RequestPost, url string, bot bot.MattermostBot) error {
	bytesRepresentation, err := json.Marshal(post)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url+"/api/v4/posts", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return err
	}
	var bearer = "Bearer " + bot.Token
	req.Header.Add("Authorization", bearer)
	resp, err := (*http.Client)(client).Do(req)
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("MattermostBot create post: %s", bytesResp)
	return nil
}
