package chat

// user http request body
type ChatRequest struct {
	ProjectName string `json:"projectName"`
	TaskID      string `json:"taskID"`
	MessageID   string `json:"messageID"`
	ThreadID    string `json:"threadID"`
	Content     string `json:"content"`
	Type        string `json:"type"`
}
