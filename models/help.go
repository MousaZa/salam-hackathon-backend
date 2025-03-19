package models

type HelpRequest struct {
	Language  string `json:"language"`
	FrameWork string `json:"framework"`
	Project   string `json:"project"`
	Task      string `json:"task"`
}

func (h *HelpRequest) ToPrompt() string {
	return `
	أقوم بعمل مشروع بعنوان ` + h.Project + ` وأواجه مشكلة في عمل ` + h.Task + ` في ` + h.FrameWork + ` بلغة البرمجة ` + h.Language + ` هل يمكنك مساعدتي.
	`
}
