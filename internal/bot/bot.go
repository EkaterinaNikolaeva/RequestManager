package bot

import (
	"os"
)

type MattermostBot struct {
	Token string
}

func LoadMattermostBot() MattermostBot {
	var settings MattermostBot
	settings.Token = os.Getenv("MATTERMOST_TOKEN")
	return settings
}
