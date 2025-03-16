package models

type Task struct {
	Id          string `json:"id"`
	Order       int    `json:"order"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	ProjectId   string `json:"projectId"`
}

func (t *Task) ToMap() map[string]interface{} {
	return map[string]interface{}{
		// "id":          t.Id,
		"order":       t.Order,
		"projectId":   t.ProjectId,
		"title":       t.Title,
		"description": t.Description,
		"completed":   t.Completed,
	}
}
