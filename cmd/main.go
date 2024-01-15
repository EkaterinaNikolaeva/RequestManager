package main

import (
	"log"
	"net/http"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/mattermostmessages"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/server"
)

func main() {
	mattermostBot := bot.LoadMattermostBot()
	client := mattermostmessages.HttpClient(http.Client{})
	err := client.SendMessage(mattermostmessages.Message{
		Message:   "abacaba",
		ChannelId: "9gs6do7otff9fmgcrktnk9opra",
	}, "http://localhost:8065", mattermostBot)
	if err != nil {
		log.Printf("%q\n", err)
	}
	server.MakeServer(mattermostBot, "http://localhost:8065")

}
