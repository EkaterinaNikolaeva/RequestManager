package config

import (
	"errors"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MattermostToken     string `yaml:"env_mattermost_token"`
	MattermostHttp      string `yaml:"mattermost_http"`
	MattermostWebsocket string `yaml:"mattermost_websocket,omitempty"`
	TeamName            string `yaml:"team_name"`
	MessagesPattern     string `yaml:"messages_pattern"`
	MessageReply        string `yaml:"message_reply"`
	JiraBotUsername     string `yaml:"env_jira_bot_username"`
	JiraBotPassword     string `yaml:"env_jira_bot_password"`
	JiraProject         string `yaml:"jira_project"`
	JiraIssueType       string `yaml:"jira_issue_type"`
	JiraBaseUrl         string `yaml:"jira_base_url"`
}

func (c *Config) getEnvVars() {
	c.MattermostToken = os.Getenv(c.MattermostToken)
	c.JiraBotUsername = os.Getenv(c.JiraBotUsername)
	c.JiraBotPassword = os.Getenv(c.JiraBotPassword)
}

var validHttp = regexp.MustCompile(`http[s]?://.*`)
var validWs = regexp.MustCompile(`ws://.*`)

func (c *Config) validateConfig() error {
	if !validWs.MatchString(c.MattermostWebsocket) {
		return errors.New("incorrect websocket server")
	}
	if !validHttp.MatchString(c.MattermostHttp) || !validHttp.MatchString(c.JiraBaseUrl) {
		return errors.New("incorrect http server")
	}
	if c.MattermostToken == "" {
		return errors.New("incorrect mattermost token")
	}
	return nil
}

func LoadConfig(fileName string) (Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	config.getEnvVars()
	err = config.validateConfig()
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
