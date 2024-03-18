package jiratasks

type JiraTaskCreationRequest struct {
	Fields JiraTaskCreationFields `json:"fields"`
}

type JiraTaskCreationFields struct {
	Project     JiraTaskCreationProject   `json:"project"`
	Summary     string                    `json:"summary"`
	Description string                    `json:"description"`
	IssueType   JiraTaskCreationIssueType `json:"issuetype"`
}

type JiraTaskCreationProject struct {
	Key string `json:"key"`
}

type JiraTaskCreationIssueType struct {
	Name string `json:"name"`
}

type JiraTaskCreationResponse struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}
