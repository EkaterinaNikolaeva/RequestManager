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

func LoadConfig(fileName string) (Config, error) {
	log.Println(fileName)
	file, err := os.Open(fileName)
	if err != nil {
		return Config{}, err
	}
	data := make([]byte, 1024)
	size, err := file.Read(data)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = yaml.Unmarshal(data[:size], &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
