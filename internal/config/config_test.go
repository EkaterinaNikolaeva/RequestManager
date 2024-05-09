package config

import (
	"io/ioutil"
	"os"
	"strconv"
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

func makeConfigDataMattermost(envMattermostToken string, mattermostHttp string, mattermostWebsocket string, teamName string, jiraHttp string) string {
	configData := "messenger: " + "mattermost" + "\n"
	configData += "env_mattermost_token: " + envMattermostToken + "\n"
	configData += "mattermost_http: " + mattermostHttp + "\n"
	configData += "mattermost_websocket: " + mattermostWebsocket + "\n"
	configData += "mettermost_team_name: " + teamName + "\n"
	configData += "jira_base_url: " + jiraHttp + "\n"
	return configData
}

func TestLoadConfig(t *testing.T) {
	envMattermostToken := "token"
	t.Setenv(envMattermostToken, envMattermostToken)
	mattermostHttp := "http://localhost:0000"
	jiraHttp := "http://localhost:0000"
	mattermostWebsocket := "ws://localhost:0000"
	teamName := "team"
	configData := makeConfigDataMattermost(envMattermostToken, mattermostHttp, mattermostWebsocket, teamName, jiraHttp)
	configFile, err := ioutil.TempFile("", "config-*.txt")
	if err == nil {
		configFile.Write([]byte(configData))
		configFile.Close()
		config, err := LoadConfig(configFile.Name())
		assert.Equal(t, nil, err)
		assert.Equal(t, config.MattermostToken, envMattermostToken)
		assert.Equal(t, config.MattermostHttp, mattermostHttp)
		assert.Equal(t, config.MattermostWebsocket, mattermostWebsocket)
		assert.Equal(t, config.MattermostTeamName, teamName)
		os.Remove(configFile.Name())
	}

}

func TestIncorrectConfig(t *testing.T) {
	configFile, err := ioutil.TempFile("", "config-*.txt")
	fileName := configFile.Name()
	if err == nil {
		configData := makeConfigDataMattermost("token", "https://localhost:0000", "ws://localhost:0000", "team", "https://localhost:0000")
		configFile.WriteString(configData)
		configFile.Close()
		_, err := LoadConfig(configFile.Name())
		assert.NotEqual(t, nil, err)
		assert.EqualError(t, err, "incorrect mattermost token")
		configFile.Close()

		configFile, _ = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 755)
		configData = makeConfigDataMattermost("token", "httpss://localhost:0000", "ws://localhost:0000", "team", "http://localhost:0000")
		configFile.Write([]byte(configData))
		configFile.Close()
		_, err = LoadConfig(configFile.Name())
		assert.NotEqual(t, nil, err)
		assert.EqualError(t, err, "incorrect http server")
		configFile.Close()

		configFile, _ = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 755)
		configData = makeConfigDataMattermost("token", "https://localhost:0000", "ws:/localhost:0000", "team", "https://localhost:0000")
		configFile.Write([]byte(configData))
		configFile.Close()
		_, err = LoadConfig(configFile.Name())
		assert.NotEqual(t, nil, err)
		assert.EqualError(t, err, "incorrect websocket server")
		configFile.Close()

		os.Remove(configFile.Name())
	}

}

func makeConfigDataWithStorage(envMattermostToken string, mattermostHttp string, mattermostWebsocket string, teamName string,
	jiraHttp string, enableMsgThreting bool, storageType string) string {
	configData := makeConfigDataMattermost(envMattermostToken, mattermostHttp, mattermostWebsocket, teamName, jiraHttp)
	configData += "enable_msg_threating: " + strconv.FormatBool(enableMsgThreting) + "\n"
	configData += "storage_type: " + storageType + "\n"
	return configData
}

func TestStorage(t *testing.T) {
	configFile, err := ioutil.TempFile("", "config-*.txt")
	t.Setenv("token", "token")
	fileName := configFile.Name()
	if err == nil {
		configData := makeConfigDataWithStorage("token", "http://localhost:0000", "ws://localhost:0000", "team", "http://localhost:0000", true, "in_memory")
		configFile.WriteString(configData)
		configFile.Close()
		config, err := LoadConfig(configFile.Name())
		assert.Equal(t, nil, err)
		assert.Equal(t, IN_MEMORY, config.StorageType)

		configFile, _ = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 755)
		configData = makeConfigDataWithStorage("token", "http://localhost:0000", "ws://localhost:0000", "team", "http://localhost:0000", true, "postgres")
		configFile.WriteString(configData)
		configFile.Close()
		config, err = LoadConfig(configFile.Name())
		assert.Equal(t, nil, err)
		assert.Equal(t, POSTGRES, config.StorageType)

		configFile, _ = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 755)
		configData = makeConfigDataWithStorage("token", "http://localhost:0000", "ws://localhost:0000", "team", "http://localhost:0000", true, "other")
		configFile.WriteString(configData)
		configFile.Close()
		config, err = LoadConfig(configFile.Name())
		assert.NotEqual(t, nil, err)
		assert.Equal(t, "incorrect storage type", err.Error())
		os.Remove(configFile.Name())
	}

}
