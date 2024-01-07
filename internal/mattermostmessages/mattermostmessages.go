package mattermostmessages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
)

type Message struct {
	Text      string `json:"message"`
	ChannelId string `json:"channel_id"`
}

func SendMessage(msg Message, url string, bot bot.MattermostBot) error {
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
	client := &http.Client{}
	resp, err := client.Do(req)
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("MattermostBot send message: " + string(bytesResp))
	return nil
}
