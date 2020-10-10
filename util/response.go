package util

import (
	"bytes"
	"encoding/json"
	"github.com/devingen/api-core/dvnruntime"
	"github.com/devingen/api-core/model"
	"net/http"
)

func BuildResponse(statusCode int, data interface{}, err error) (dvnruntime.Response, error) {

	if err != nil {
		// return dvn error in the response body
		switch castedError := err.(type) {
		case model.DVNError:
			statusCode = castedError.StatusCode
			data = castedError
		case *model.DVNError:
			statusCode = castedError.StatusCode
			data = castedError
		default:
			statusCode = 418
			data = model.NewStatusError(http.StatusInternalServerError)
		}
	}

	body, err := json.Marshal(data)
	if err != nil {
		return dvnruntime.Response{StatusCode: 500}, nil
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	return dvnruntime.Response{
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
