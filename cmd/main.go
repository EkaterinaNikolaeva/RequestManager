package main

import (
	"log"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/bot"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/config"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/server"
)

func main() {
	mattermostBot := bot.LoadMattermostBot()
	config := config.LoadConfig("../configs/config.json")
	log.Printf("%s %s %s %s", config.EnvMattermostToken, config.MattermostHttp, config.MattermostWebsocket, config.TeamName)
	server.MakeServer(mattermostBot, "http://localhost:8065")
}
