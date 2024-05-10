package config

import (
	"errors"
	"html/template"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type StorageType string
type Messenger string

const (
	IN_MEMORY StorageType = "in_memory"
	POSTGRES  StorageType = "postgres"
)

const (
	MATTERMOST Messenger = "mattermost"
	ROCKETCHAT Messenger = "rocketchat"
)

type Config struct {
	Messenger Messenger `yaml:"messenger"`

	MattermostToken         string `yaml:"env_mattermost_token"`
	MattermostHttp          string `yaml:"mattermost_http"`
	MattermostWebsocket     string `yaml:"mattermost_websocket,omitempty"`
	MattermostTeamName      string `yaml:"mettermost_team_name"`
	MessagesPattern         string `yaml:"messages_pattern"`
	MessageReply            string `yaml:"message_reply"`
	JiraBotUsername         string `yaml:"env_jira_bot_username"`
	JiraBotPassword         string `yaml:"env_jira_bot_password"`
	JiraProject             string `yaml:"jira_project"`
	JiraIssueType           string `yaml:"jira_issue_type"`
	JiraBaseUrl             string `yaml:"jira_base_url"`
	MessagesPatternTemplate *template.Template
	EnableMsgThreating      bool        `yaml:"enable_msg_threating"`
	StorageType             StorageType `yaml:"storage_type"`
	PostgresLogin           string      `yaml:"env_postgres_login"`
	PostgresPassword        string      `yaml:"env_postgres_password"`
	PostgresHost            string      `yaml:"postgres_host"`
	PostgresPort            string      `yaml:"postgres_port"`
	PostgresName            string      `yaml:"postgres_name"`
	PostgresTableName       string      `yaml:"postgres_table_name"`
	RocketchatHost          string      `yaml:"rocketchat_host"`
	RocketchatToken         string      `yaml:"env_rocketchat_token"`
	RocketchatId            string      `yaml:"env_rocketchat_id"`
	RocketchatHttp          string      `yaml:"rocketchat_http"`
}

func (c *Config) getEnvVars() {
	c.MattermostToken = os.Getenv(c.MattermostToken)
	c.JiraBotUsername = os.Getenv(c.JiraBotUsername)
	c.JiraBotPassword = os.Getenv(c.JiraBotPassword)
	c.RocketchatToken = os.Getenv(c.RocketchatToken)
	c.RocketchatId = os.Getenv(c.RocketchatId)
	if c.PostgresLogin != "" {
		c.PostgresLogin = os.Getenv(c.PostgresLogin)
		c.PostgresPassword = os.Getenv(c.PostgresPassword)
	}
}

func (c *Config) compileTemplates() error {
	tmpl, err := template.New("test").Parse(c.MessageReply)
	if err != nil {
		return err
	}
	c.MessagesPatternTemplate = tmpl
	return nil
}

var validHttp = regexp.MustCompile(`http[s]?://.*`)
var validWs = regexp.MustCompile(`ws://.*`)

func (c *Config) validateConfig() error {
	if c.Messenger == MATTERMOST && !validWs.MatchString(c.MattermostWebsocket) {
		return errors.New("incorrect websocket server")
	}
	if c.Messenger == MATTERMOST && !validHttp.MatchString(c.MattermostHttp) || !validHttp.MatchString(c.JiraBaseUrl) {
		return errors.New("incorrect http server")
	}
	if c.MattermostToken == "" {
		return errors.New("incorrect mattermost token")
	}
	if c.StorageType != POSTGRES && c.StorageType != IN_MEMORY && c.EnableMsgThreating {
		return errors.New("incorrect storage type")
	}
	if c.Messenger != MATTERMOST && c.Messenger != ROCKETCHAT {
		return errors.New("incorrect messenger")
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
	config.compileTemplates()
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
