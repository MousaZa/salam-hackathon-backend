package models

type LearningRequest struct {
	Language  string `json:"language"`
	Level     string `json:"level"`
	FrameWork string `json:"framework"`
	Goal      string `json:"goal"`
}

type LearningResponse struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Projects    []Project `json:"projects"`
}

func (l *LearningResponse) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"title":       l.Title,
		"description": l.Description,
		"projects":    l.Projects,
	}
}
