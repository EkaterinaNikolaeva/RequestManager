package bot

import (
	"github.com/EkaterinaNikolaeva/RequestManager/internal/config"
)

type MattermostBot struct {
	Token               string
	TeamName            string
	MattermostHttp      string
	MattermostWebsocket string
}

func NewMattermostBot(config config.Config) MattermostBot {
	var bot MattermostBot
	bot.Token = config.MattermostToken
	bot.TeamName = config.TeamName
	bot.MattermostHttp = config.MattermostHttp
	bot.MattermostWebsocket = config.MattermostWebsocket
	return bot
}
