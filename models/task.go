package models

import "fmt"

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

func (t *Task) String() string {
	return fmt.Sprintf("Title: %s, Description: %s, Completed: %v\n", t.Title, t.Description, t.Completed)
}
