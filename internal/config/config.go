package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	EnvMattermostToken  string `json:"env_mattermost_token"`
	MattermostHttp      string `json:"mattermost_http"`
	MattermostWebsocket string `json:"mattermost_websocket,omitempty"`
	TeamName            string `json:"team_name"`
}

func LoadConfig(fileName string) Config {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error when opening config file: %q", err)
	}
	data := make([]byte, 1024)
	size, err := file.Read(data)
	if err != nil {
		log.Fatalf("Error when reading config file: %q", err)
	}
	var config Config
	err = json.Unmarshal(data[:size], &config)
	if err != nil {
		log.Fatalf("Error when encode config file %q", err)
	}
	return config
}
