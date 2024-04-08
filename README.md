# RequestManager

## Configuration 
An example of the config is in the file `configs/config.yml`

* `env_mattermost_token` - the name of the environment variable containing the mattermost bot token
* `env_jira_bot_username` и `env_jira_bot_password` - the names of the environment variables containing the username and password of the bot's account in jira
* `jira_project`, `jira_issue_type` allows you to configure the task creation parameters in jira
* `messages_pattern` - the regular expression that will determine the need to create an issue
* `message_reply` - template for generating the default response of the bot in the messenger.

## Installation
* Install `mattermost`, `jira` 
* Create bot accounts
* Save the necessary tokens to environment variables
* Set up the configuration file
* Build and run the application from the cmd directory
```
go build
./cmd "../configs/config.yml"
```