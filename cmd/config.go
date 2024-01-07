package main

import (
	"os"
)

// type config struct {
// 	mattermostToken string
// }

func loadConfig() config {
	var settings config
	settings.MattermostToken = os.Getenv("MATTERMOST_TOKEN")
	return settings
}
