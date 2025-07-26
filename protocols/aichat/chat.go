package aichat

type ChatRequest struct {
	ThreadID string `json:"threadID"`
	Content  string `json:"content"`
	Type     string `json:"type"`
}
