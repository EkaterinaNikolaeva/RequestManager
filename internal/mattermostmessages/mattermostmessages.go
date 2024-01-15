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

type Message struct {
	Id            string                 `json:"id,omitempty"`
	CreateAt      int                    `json:"create_at,omitempty"`
	UpdateAt      int                    `json:"update_at,omitempty"`
	EditAt        int                    `json:"edit_at,omitempty"`
	DeleteAt      int                    `json:"delete_at,omitempty"`
	IsPinned      bool                   `json:"is_pinned,omitempty"`
	UserId        string                 `json:"user_id,omitempty"`
	ChannelId     string                 `json:"channel_id"`
	RootId        string                 `json:"root_id,omitempty"`
	OriginalId    string                 `json:"original_id,omitempty"`
	Message       string                 `json:"message"`
	Props         map[string]bool        `json:"props,omitempty"`
	Hashtags      string                 `json:"hashtags,omitempty"`
	PendingPostId string                 `json:"pending_post_id,omitempty"`
	ReplyCount    int                    `json:"reply_count,omitempty"`
	LastReplyAt   int                    `json:"last_reply_at,omitempty"`
	Participants  map[string]interface{} `json:"participants,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

var validMsg = regexp.MustCompile(`.*jira!.*`)

func checkPatternInMessage(msg string) bool {
	return validMsg.MatchString(strings.ToLower(msg))
}

func CheckMessageForJiraRequest(bytes string) {
	var msg Message
	err := json.Unmarshal([]byte(bytes), &msg)
	if err != nil {
		log.Printf("Error when encode message %q", err)
		return
	}
	if checkPatternInMessage(msg.Message) {
		log.Printf("Message: %s, make issue!", msg.Message)
	}
}

type HttpClient http.Client

func (client *HttpClient) SendMessage(msg Message, url string, bot bot.MattermostBot) error {
	bytesRepresentation, err := json.Marshal(msg)
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
	log.Println("MattermostBot send message: " + string(bytesResp))
	return nil
}
