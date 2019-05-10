package views

import ()

type SystemMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type ErrorMessage struct {
	Message SystemMessage `json:"system_message"`
	Errors  []string      `json:"errors"`
}

type MainPagination struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalEntries int `json:"total_entries"`
}

func SetSystemMessage(mType string, content string) SystemMessage {
	return SystemMessage{Type: mType, Content: content}
}

func SetErrorMessage(mType string, content string, mErrors []string) ErrorMessage {
	return ErrorMessage{Message: SetSystemMessage(mType, content), Errors: mErrors}
}

func SetNotFoundErrorMessage() ErrorMessage {
	var errors []string
	return SetErrorMessage("alert", "404 not found", errors)
}
