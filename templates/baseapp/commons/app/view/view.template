package view

import ()

type DefaultMessage struct {
	Message SystemMessage `json:"system_message"`
}

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

func SetDefaultMessage(mType string, content string) DefaultMessage {
	return DefaultMessage{Message: SetSystemMessage(mType, content)}
}

func SetErrorMessage(mType string, content string, errs []error) ErrorMessage {
	var stringErrors []string

	for _, value := range errs {
		stringErrors = append(stringErrors, value.Error())
	}

	return ErrorMessage{Message: SetSystemMessage(mType, content), Errors: stringErrors}
}

func SetNotFoundErrorMessage() ErrorMessage {
	var errs []error
	return SetErrorMessage("alert", "404 not found", errs)
}
