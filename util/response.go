package util

import (
	"bytes"
	"encoding/json"
	"github.com/devingen/api-core/dto"
)

func BuildResponse(statusCode int, data interface{}, err error) (dto.Response, error) {

	if err != nil {
		return dto.Response{StatusCode: 500}, nil
	}

	body, err := json.Marshal(data)
	if err != nil {
		return dto.Response{StatusCode: 500}, nil
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	return dto.Response{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
	}, nil
}
