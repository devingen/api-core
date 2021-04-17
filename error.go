package core

import "net/http"

type DVNError struct {
	Message    string `json:"error,omitempty"`
	ErrorCode  int    `json:"errorCode,omitempty"`
	StatusCode int    `json:"-"`
}

func (fe DVNError) Error() string {
	if fe.Message == "" {
		return http.StatusText(fe.StatusCode)
	}
	return fe.Message
}

func NewStatusError(statusCode int) *DVNError {
	return &DVNError{
		StatusCode: statusCode,
		Message:    http.StatusText(statusCode),
	}
}

func NewError(statusCode int, errorMessage string) *DVNError {
	return &DVNError{
		Message:    errorMessage,
		StatusCode: statusCode,
	}
}

func NewErrorWithCode(statusCode, errorCode int, errorMessage string) *DVNError {
	return &DVNError{
		Message:    errorMessage,
		ErrorCode:  errorCode,
		StatusCode: statusCode,
	}
}
