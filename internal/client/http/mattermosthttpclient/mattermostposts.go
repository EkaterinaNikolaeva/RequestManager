package mattermosthttpclient

type RequestPost struct {
	ChannelId string                 `json:"channel_id"`
	Message   string                 `json:"message"`
	RootId    string                 `json:"root_id,omitempty"`
	FileIds   []string               `json:"file_ids,omitempty"`
	Props     interface{}            `json:"props,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type ResponsePost struct {
	Id            string                 `json:"id,omitempty"`
	CreateAt      int                    `json:"create_at,omitempty"`
	UpdateAt      int                    `json:"update_at,omitempty"`
	DeleteAt      int                    `json:"delete_at,omitempty"`
	EditAt        int                    `json:"edit_at,omitempty"`
	UserId        string                 `json:"user_id,omitempty"`
	ChannelId     string                 `json:"channel_id,omitempty"`
	RootId        string                 `json:"root_id,omitempty"`
	OriginalId    string                 `json:"original_id,omitempty"`
	Message       string                 `json:"message,omitempty"`
	Type          string                 `json:"type,omitempty"`
	Props         map[string]interface{} `json:"props,omitempty"`
	Hashtag       string                 `json:"hashtag,omitempty"`
	PendingPostId string                 `json:"pending_post_id,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}
