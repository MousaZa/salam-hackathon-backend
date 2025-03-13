package models

import "fmt"

type LearningRequest struct {
	Language  string `json:"language"`
	Level     string `json:"level"`
	FrameWork string `json:"framework"`
	Goal      string `json:"goal"`
}

func (l *LearningRequest) ToPrompt() string {
	return fmt.Sprintf(`response in Arabic. I want to learn %s in %s programming language, I am a %s and my goal is to %s. List 5 projects to work on to develop my skills using this JSON schema:
        LearningResponse: { 'title' : string, 'description' : string ,'projects' : [{'title': string , 'description': string, 'tasks': [{'title': string, 'description': string, 'completed': bool}]}]}       
		Return:   LearningResponse
	           `, l.FrameWork, l.Language, l.Level, l.Goal)
}

type LearningResponse struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Projects    []Project `json:"projects"`
}

func (l *LearningResponse) String() string {
	return fmt.Sprintf("Title: %s, Description: %s,", l.Title, l.Description)
}

func (l *LearningResponse) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"title":       l.Title,
		"description": l.Description,
		"projects":    l.Projects,
	}
}
