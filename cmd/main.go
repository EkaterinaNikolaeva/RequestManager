package main

import (
	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/server"
)

func main() {
	mattermostBot := bot.LoadMattermostBot()
	server.MakeServer(mattermostBot, "http://localhost:8065")
}
