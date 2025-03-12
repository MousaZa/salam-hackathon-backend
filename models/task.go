package models

type Task struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func (t *Task) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"title":       t.Title,
		"description": t.Description,
		"completed":   t.Completed,
	}
}
