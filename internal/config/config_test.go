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
		assert.EqualError(t, err, "open "+fileDeleted.Name()+": no such file or directory")
	}
}

func makeConfigData(envMattermostToken string, mattermostHttp string, mattermostWebsocket string, teamName string) string {
	configData := "env_mattermost_token: " + envMattermostToken + "\n"
	configData += "mattermost_http: " + mattermostHttp + "\n"
	configData += "mattermost_websocket: " + mattermostWebsocket + "\n"
	configData += "team_name: " + teamName + "\n"
	return configData
}

func TestLoadConfig(t *testing.T) {
	envMattermostToken := "token"
	mattermostHttp := "http://localhost:0000"
	mattermostWebsocket := "ws://localhost:0000"
	teamName := "team"
	configData := makeConfigData(envMattermostToken, mattermostHttp, mattermostWebsocket, teamName)
	configFile, err := ioutil.TempFile("", "config-*.txt")
	if err == nil {
		configFile.Write([]byte(configData))
		configFile.Close()
		config, err := LoadConfig(configFile.Name())
		assert.Equal(t, nil, err)
		assert.Equal(t, config.MattermostToken, envMattermostToken)
		assert.Equal(t, config.MattermostHttp, mattermostHttp)
		assert.Equal(t, config.MattermostWebsocket, mattermostWebsocket)
		assert.Equal(t, config.TeamName, teamName)
		os.Remove(configFile.Name())
	}

}

func TestIncorrectConfig(t *testing.T) {
	configFile, err := ioutil.TempFile("", "config-*.txt")
	fileName := configFile.Name()
	if err == nil {
		configData := makeConfigData("token", "https://localhost:0000", "ws://localhost:0000", "team")
		configFile.WriteString(configData)
		configFile.Close()
		_, err := LoadConfig(configFile.Name())
		assert.Equal(t, nil, err)
		configFile.Close()

		configFile, _ = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 755)
		configData = makeConfigData("token", "httpss://localhost:0000", "ws://localhost:0000", "team")
		configFile.Write([]byte(configData))
		configFile.Close()
		_, err = LoadConfig(configFile.Name())
		assert.NotEqual(t, nil, err)
		assert.EqualError(t, err, "incorrect http server")
		configFile.Close()

		configFile, _ = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 755)
		configData = makeConfigData("token", "https://localhost:0000", "ws:/localhost:0000", "team")
		configFile.Write([]byte(configData))
		configFile.Close()
		_, err = LoadConfig(configFile.Name())
		assert.NotEqual(t, nil, err)
		assert.EqualError(t, err, "incorrect websocket server")
		configFile.Close()

		os.Remove(configFile.Name())
	}

}
