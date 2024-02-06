package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	EnvMattermostToken  string `yaml:"env_mattermost_token"`
	MattermostHttp      string `yaml:"mattermost_http"`
	MattermostWebsocket string `yaml:"mattermost_websocket,omitempty"`
	TeamName            string `yaml:"team_name"`
}

func LoadConfig(fileName string) Config {
	log.Println(fileName)
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
	err = yaml.Unmarshal(data[:size], &config)
	if err != nil {
		log.Fatalf("Error when encode config file %q", err)
	}
	return config
}
