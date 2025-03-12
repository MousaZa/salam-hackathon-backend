package models

type Project struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Tasks       []Task `json:"tasks"`
}

func (p *Project) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"title":       p.Title,
		"description": p.Description,
		"tasks":       p.Tasks,
	}
}
