package rocketchathttpclient

type RequestMessage struct {
	Message RequestMessageData `json:"message"`
}

type RequestMessageData struct {
	Rid  string `json:"rid"`
	Msg  string `json:"msg"`
	Tmid string `json:"tmid"`
}

type ResponsePost struct {
	Success string `json:"success"`
	Message string `json:"message"`
}

type RespondeMessageData struct {
	Alias     string           `json:"alias"`
	Msg       string           `json:"msg"`
	ParseUrls string           `json:"parseUrls"`
	Groupable string           `json:"groupable"`
	Ts        int64            `json:"ts"`
	U         RespondeUserData `json:"u"`
	Rid       string           `json:"rid"`
	Tmid      string           `json:"tmid"`
	UpdatedAt int64            `json:"_updatedAt"`
	Id        string           `json:"_id"`
	Success   bool             `json:"success"`
	Urls      []string         `json:"urls"`
	Mentions  []string         `json:"mentions"`
	Channels  []string         `json:"channels"`
	Md        []interface{}    `json:"md"`
}

type RespondeUserData struct {
	Id       string `json:"_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}
