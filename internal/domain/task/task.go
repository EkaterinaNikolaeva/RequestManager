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
