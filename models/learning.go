package models

import "fmt"

type LearningRequest struct {
	Language  string `json:"language"`
	Level     string `json:"level"`
	FrameWork string `json:"framework"`
	Goal      string `json:"goal"`
}

func (l *LearningRequest) ToPrompt() string {
	return fmt.Sprintf(`response in Arabic. I want to learn %s in %s programming language, I am a %s and my goal is to %s. List 5 projects to work on to develop my skills (in order depending on difficulty, make the fields no too long using this JSON schema:
        LearningResponse: { 'title' : string, 'description' : string ,'projects' : [{'title': string , 'description': string, 'order' : int, 'tasks': [{'title': string, 'order' : int, 'description': string, 'completed': bool}]}]}       
		Return:   LearningResponse
	           `, l.FrameWork, l.Language, l.Level, l.Goal)
}

type LearningResponse struct {
	Id          string    `json:"id"`
	Language    string    `json:"language"`
	Level       string    `json:"level"`
	FrameWork   string    `json:"framework"`
	Goal        string    `json:"goal"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Projects    []Project `json:"projects"`
}

func (l *LearningResponse) ToMap() map[string]interface{} {
	return map[string]interface{}{
		// "id":          l.Id,
		"language":    l.Language,
		"level":       l.Level,
		"framework":   l.FrameWork,
		"goal":        l.Goal,
		"title":       l.Title,
		"description": l.Description,
		// "projects":    l.Projects,
	}
}
