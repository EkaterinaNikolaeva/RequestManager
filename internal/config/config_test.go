package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigNoSuchFile(t *testing.T) {
	fileDeleted, err := ioutil.TempFile("", "config-*.txt")
	if err == nil {
		fileDeleted.Close()
		os.Remove(fileDeleted.Name())
		_, err := LoadConfig(fileDeleted.Name())
		assert.NotEqual(t, nil, err)
	}
}

func TestLoadConfig(t *testing.T) {
	envMattermostToken := "token"
	mattermostHttp := "http://localhost:0000"
	mattermostWebsocket := "ws://localhost:0000"
	teamName := "team"
	configData := "env_mattermost_token: " + envMattermostToken + "\n"
	configData += "mattermost_http: " + mattermostHttp + "\n"
	configData += "mattermost_websocket: " + mattermostWebsocket + "\n"
	configData += "team_name: " + teamName + "\n"
	configFile, err := ioutil.TempFile("", "config-*.txt")
	if err == nil {
		configFile.Write([]byte(configData))
		configFile.Close()
		config, err := LoadConfig(configFile.Name())
		assert.Equal(t, nil, err)
		assert.Equal(t, config.EnvMattermostToken, envMattermostToken)
		assert.Equal(t, config.MattermostHttp, mattermostHttp)
		assert.Equal(t, config.MattermostWebsocket, mattermostWebsocket)
		assert.Equal(t, config.TeamName, teamName)
		os.Remove(configFile.Name())
	}

}
