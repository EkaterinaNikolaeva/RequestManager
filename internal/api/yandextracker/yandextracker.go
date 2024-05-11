package apiyandextracker

type RequestTask struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Queue       string `json:"queue"`
	Type        string `json:"type"`
}

type ResponseTask struct {
	Self                               string             `json:"self"`
	Id                                 string             `json:"id"`
	Key                                string             `json:"key"`
	Version                            int                `json:"version"`
	Summary                            string             `json:"summary"`
	StatusStartTime                    string             `json:"statusStartTime"`
	UpdatedBy                          ResponseUser       `json:"updatedBy"`
	StatusType                         interface{}        `json:"statusType"`
	Description                        string             `json:"description"`
	Type                               ResponseParamsTask `json:"type"`
	Priority                           ResponseParamsTask `json:"priority"`
	CreatedAt                          string             `json:"createdAt"`
	CreatedBy                          ResponseUser       `json:"createdBy"`
	CommentWithoutExternalMessageCount int                `json:"commentWithoutExternalMessageCount"`
	Votes                              int                `json:"votes"`
	Queue                              ResponseParamsTask `json:"queue"`
	Status                             ResponseParamsTask `json:"status"`
	Favorite                           bool               `json:"favorite"`
}

type ResponseUser struct {
	Self        string `json:"self"`
	Id          string `json:"id"`
	Display     string `json:"display"`
	CloudUid    string `json:"cloudUid"`
	PassportUid int    `json:"passportUid"`
}

type ResponseParamsTask struct {
	Self    string `json:"self"`
	Id      string `json:"id"`
	Key     string `json:"key"`
	Display string `json:"display"`
}

type RequestComment struct {
	Text string `json:"text"`
}

type ResponseComment struct {
	Self              string                  `json:"self"`
	Id                string                  `json:"id"`
	LongId            string                  `json:"longId"`
	Text              string                  `json:"text"`
	CreateBody        ResponseParamsComment   `json:"createBody"`
	UpdateBody        ResponseParamsComment   `json:"updateBody"`
	CreatedAt         string                  `json:"createdAt"`
	UpdatedAt         string                  `json:"updatedAt"`
	Summonees         []ResponseParamsComment `json:"summonees"`
	MaillistSummonees []ResponseParamsComment `json:"maillistSummonees"`
	Version           int                     `json:"version"`
	Type              string                  `json:"type"`
	Transport         string                  `json:"transport"`
}

type ResponseParamsComment struct {
	Self    string `json:"self"`
	Id      string `json:"id"`
	Display string `json:"display"`
}
