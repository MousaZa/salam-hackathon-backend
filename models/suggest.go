package models

import "fmt"

type SuggestRequest struct {
	Goal       string `json:"goal"`
	Level      string `json:"level"`
	Preference string `json:"preference"`
}

func (sr *SuggestRequest) ToPrompt() string {
	return fmt.Sprintf(`
		response in Arabic. My goal is to %s, I am a %s with coding and I prefer to work with %s. Suggest me a programming language and a framework to work on.
	`, sr.Goal, sr.Level, sr.Preference)
}
