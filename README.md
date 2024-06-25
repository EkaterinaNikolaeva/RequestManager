# RequestManager

* [About](#About)
* [Installation](#Installation)
* [Configuration](#Configuration)
* [Screenshots](#Screenshots)

## About
A service that automates the process of setting tasks in the task tracker by analyzing messages in the messenger. Integration with Mattermost, Rocket Chat, Jira and Yandex Tracker is currently supported. It is possible to track new messages in the threads for which tasks have been created, with further addition of the message text in the comments to the task.

## Installation
* Install the necessary messenger (`Mattermost` or `Rocket.Chat`) and task tracker (`Jira` or `Yandex Tracker`)
* Create bot accounts
* Save the necessary tokens to environment variables
* Set up the configuration file
* Add the name of the file containing the config to the environment variable `CONFIG`
* Build and run
```
make run
```

## Configuration 
An example of the config is in the file `configs/config.yml`
* `messenger` - the messenger used. `rocketchat` or `mattermost`
* `task_tracker` - the task tracker used. `jira` or `yandex_tracker`
* `enable_msg_threating` - the mode for tracking new messages for which a task has been created, and creating comments on the task with text. `true` or `false`
* `storage_type` - `postgres` or `in_memory`. It is necessary if `enable_msg_threating == true`, 
* `message_reply` - template for generating the default response of the bot in the messenger
* `messages_pattern` - the regular expression that will determine the need to create an issue
* `task_name` - template for generating the default task name

If `storage_type == postgres`, then the following parameters must also be specified:

* `env_postgres_login` - the name of the environment variable containing postgres login
* `env_postgres_password` - the name of the environment variable containing postgres password
* `postgres_host, postgres_port, postgres_name, postgres_table_name` - the hostname, port, name and the tablename of the postgres database server

Next, you can specify only the parameters necessary for a specific messenger and task tracker

Jira:
* `env_jira_bot_username` и `env_jira_bot_password` - the names of the environment variables containing the username and password of the bot's account in jira
* `jira_project`, `jira_issue_type` allows you to configure the task creation parameters in jira
* `jira_base_url` - base url for connecting to a jira server

Mattermost:
* `env_mattermost_token` - the name of the environment variable containing the mattermost bot token
* `mattermost_http` - HTTP or HTTPS endpoint for connecting to a mattermost server
* `mattermost_websocket` - WebSocket endpoint for connecting to a mattermost server
* `mettermost_team_name` - name of the mattermost team

Rocket.Chat:
* `mattermost_http` - HTTP or HTTPS endpoint for connecting to a rocketchat server
* `mattermost_host` - base url for connecting to a rocketchat server
* `env_rocketchat_token` - the name of the environment variable containing the rocketchat bot token
* `env_rocketchat_id` - the name of the environment variable containing the rocketchat bot id

Yandex Tracker:

* `yandex_tracker_host` - API endpoint for connecting to the yandex tracker service
* `yandex_tracker_base_url` - base url for connecting to the yandex tracker service
* `env_yandex_tracker_id_organization, env_yandex_tracker_token` - the name of the environment variable containing id organization, token of the bot
* `yandex_tracker_type_organization, yandex_tracker_token_type` - the name of the environment variable containing type organization, token type of the bot
* `yandex_tracker_queue`, `yandex_tracker_task_type`: allows you to configure the task creation parameters in yandex tracker

## Screenshots
Make issue in Yansex Tracker from Rocket.Chat:

<kbd>
<a href="https://drive.google.com/uc?export=view&id=1FafBEh-SpztnF7gqo5aLkAa4zww1CthI"><img src="https://drive.google.com/uc?export=view&id=1FafBEh-SpztnF7gqo5aLkAa4zww1CthI" style="width: 500px; max-width: 100%; height: auto" title="Make issue in Yansex Tracker from Rocket.Chat" />
</kbd>

Task in Yandex Tracker:

<kbd>
<a href="https://drive.google.com/uc?export=view&id=1F6v1PnJXgkXhpt5oAzJQqUE4daMvnVtK"><img src="https://drive.google.com/uc?export=view&id=1F6v1PnJXgkXhpt5oAzJQqUE4daMvnVtK" style="width: 500px; max-width: 100%; height: auto" title="Task in Yandex Tracker" />
</kbd>

Make issue in Jira from Mattermost with comment:

<kbd>
<a href="https://drive.google.com/uc?export=view&id=1VsM-VIjz-My8QM0LtWmzmdSL1sjjqRvk"><img src="https://drive.google.com/uc?export=view&id=1VsM-VIjz-My8QM0LtWmzmdSL1sjjqRvk" style="width: 500px; max-width: 100%; height: auto" title="Make issue in Jira from Mattermost with comment" />
</kbd>

Task in Jira with comment:

<kbd>
<a href="https://drive.google.com/uc?export=view&id=1PqvEMavwIBajgcs2vDDRSEX7F16fbbNV"><img src="https://drive.google.com/uc?export=view&id=1PqvEMavwIBajgcs2vDDRSEX7F16fbbNV" style="width: 500px; max-width: 100%; height: auto" title="Task in Jira with comment" />
</kbd>
