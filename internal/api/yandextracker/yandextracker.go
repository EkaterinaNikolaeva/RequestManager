package apiyandextracker

type RequestTask struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Queue       string `json:"queue"`
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
