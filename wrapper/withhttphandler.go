package wrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/util"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// WithHTTPHandler wraps the controller with the HTTP Handler func.
func WithHTTPHandler(ctx context.Context, f core.Controller) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// convert HTTP request to our custom request
		req := adaptHTTPRequest(r)

		// execute function
		result, status, err := f(ctx, req)

		// convert response to our custom response
		response, err := buildHTTPResponse(status, result, err)

		// write response data
		returnHTTPResponse(w, response, err)
	}
}

func adaptHTTPRequest(req *http.Request) core.Request {
	body, _ := ioutil.ReadAll(req.Body)  // TODO handle error
	ip, _ := util.GetClientIPHelper(req) // TODO handle error

	return core.Request{
		RemoteAddr:            req.RemoteAddr,
		Proto:                 req.Proto,
		Resource:              req.RequestURI,
		Path:                  req.URL.Path,
		HTTPMethod:            req.Method,
		Headers:               convertHeaders(req.Header),
		QueryStringParameters: convertQueryParams(req.URL.Query()),
		PathParameters:        mux.Vars(req),
		StageVariables:        nil,
		RequestContext:        core.ProxyRequestContext{},
		Body:                  string(body),
		IsBase64Encoded:       false,
		IP:                    ip,
	}
}

func convertHeaders(header http.Header) map[string]string {
	headers := map[string]string{}
	for k, v := range header {
		headers[k] = v[0]
	}
	return headers
}

func convertQueryParams(values url.Values) map[string]string {
	params := map[string]string{}
	for k, v := range values {
		params[k] = v[0]
	}
	return params
}

func buildHTTPResponse(statusCode int, data interface{}, err error) (core.Response, error) {

	if err != nil {
		// return dvn error in the response body
		switch castedError := err.(type) {
		case core.DVNError:
			statusCode = castedError.StatusCode
			data = castedError
		case *core.DVNError:
			statusCode = castedError.StatusCode
			data = castedError
		default:
			statusCode = http.StatusInternalServerError
			data = core.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	resp := core.Response{
		StatusCode:      statusCode,
		IsBase64Encoded: false,

		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
	}

	if data != nil {
		body, jsonError := json.Marshal(data)
		if jsonError != nil {
			return core.Response{StatusCode: 500}, nil
		}

		var buf bytes.Buffer
		json.HTMLEscape(&buf, body)
		resp.Body = buf.String()
	}
	return resp, nil
}

func returnHTTPResponse(w http.ResponseWriter, response core.Response, err error) {

	// set response headers
	if response.Headers != nil {
		for name, value := range response.Headers {
			w.Header().Set(name, value)
		}
	}

	// set response status
	if response.StatusCode != 0 {
		w.WriteHeader(response.StatusCode)
	} else if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// return response
	if response.Body != "" {
		if response.StatusCode == 204 {
			fmt.Println("--")
			fmt.Println(response.Body)
			fmt.Println("--")
		}
		_, err = w.Write([]byte(response.Body))
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}
