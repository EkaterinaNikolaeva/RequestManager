package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type config struct {
	MattermostToken string `json:"mattermost_token"`
}

type Message struct {
	Text      string `json:"message"`
	ChannelId string `json:"channel_id"`
}

func main() {
	config := loadConfig()
	data, _ := json.Marshal(config)
	var msg = Message{
		Text:      "abacaba",
		ChannelId: "9gs6do7otff9fmgcrktnk9opra",
	}
	bytesRepresentation, err := json.Marshal(msg)
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:8065/api/v4/posts", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}
	var bearer = "Bearer " + config.MattermostToken
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(bytesResp))

	fmt.Printf("%s\n", data)
}
