package mattermostmessages

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
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
	Props         interface{}            `json:"props,omitempty"`
	Hashtag       string                 `json:"hashtag,omitempty"`
	PendingPostId string                 `json:"pending_post_id,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

var validMsg = regexp.MustCompile(`.*jira!.*`)

func checkPatternInMessage(msg string) bool {
	return validMsg.MatchString(strings.ToLower(msg))
}

func (client *HttpClient) CheckMessageForJiraRequest(bytes string, mattermostBot bot.MattermostBot) {
	var post ResponsePost
	err := json.Unmarshal([]byte(bytes), &post)
	if err != nil {
		log.Printf("Error when encode message %q", err)
		return
	}
	if checkPatternInMessage(post.Message) {
		client.makePostForCreation(post, mattermostBot)
	}
}

func (client *HttpClient) GetMessage(bytes string) (string, error) {
	var post ResponsePost
	err := json.Unmarshal([]byte(bytes), &post)
	if err != nil {
		return "", err
	}
	return post.Message, nil
}
func (client *HttpClient) makePostForCreation(post ResponsePost, mattermostBot bot.MattermostBot) {
	log.Printf("Message: %s, make issue!", post.Message)
	rootId := post.Id
	if post.RootId != "" {
		rootId = post.RootId
	}
	client.CreatePost(RequestPost{
		Message:   "Create an issue. Link: ",
		ChannelId: post.ChannelId,
		RootId:    rootId,
	}, mattermostBot.MattermostHttp, mattermostBot)
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
