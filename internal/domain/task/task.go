package task

type TaskCreateRequest struct {
	Name        string
	Description string
	Type        string
	Project     string
}

type TaskCreated struct {
	Name        string
	Description string
	Type        string
	Project     string
	Link        string
}

func NewTaskCreateRequest(name string, description string, typeTask string, project string) TaskCreateRequest {
	return TaskCreateRequest{
		Name:        name,
		Description: description,
		Type:        typeTask,
		Project:     project,
	}
}
