package models

type Project struct {
	Id          string `json:"id"`
	Order       int    `json:"order"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Tasks       []Task `json:"tasks"`
	IsLocked    bool   `json:"isLocked"`
	LearningId  string `json:"learningId"`
}

func (p *Project) ToMap() map[string]interface{} {
	return map[string]interface{}{
		// "id":          p.Id,
		"title":       p.Title,
		"description": p.Description,
		"order":       p.Order,
		"learningId":  p.LearningId,
		"isLocked":    p.IsLocked,
	}
}
