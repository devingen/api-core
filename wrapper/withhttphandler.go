package wrapper

import (
	"bytes"
	"context"
	"encoding/json"
	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/util"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

// WithHTTPHandler wraps the controller with the HTTP Handler func.
func WithHTTPHandler(ctx context.Context, f core.Controller) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// convert HTTP request to our custom request
		req := adaptHTTPRequest(r)

		// execute function
		response, err := f(ctx, req)

		// convert response to our custom response
		httpResponse, err := buildHTTPResponse(response, err)

		// write response data
		returnHTTPResponse(w, httpResponse, err)
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
		QueryStringParameters: req.URL.Query(),
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

func buildHTTPResponse(response *core.Response, err error) (core.Response, error) {
	var statusCode int = 200
	var data interface{} = nil
	var headers = map[string]string{}
	var rawBody = ""

	if response != nil {
		if response.StatusCode != 0 {
			statusCode = response.StatusCode
		}
		if response.Headers != nil {
			headers = response.Headers
		}
		data = response.Body
		rawBody = response.RawBody
	}

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

	if _, ok := headers["Content-Type"]; !ok {
		// set default content type header
		headers["Content-Type"] = "application/json"
	}

	// TODO don't do this here.
	//   implement a wrapper util function and leave this to the developer's choice to add or not
	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Credentials"] = "true"

	resp := core.Response{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Headers:         headers,
		RawBody:         rawBody,
	}

	if data != nil {
		// generate the RawBody if the Body is not empty
		body, jsonError := json.Marshal(data)
		if jsonError != nil {
			return core.Response{StatusCode: 500}, nil
		}

		var buf bytes.Buffer
		json.HTMLEscape(&buf, body)
		resp.RawBody = buf.String()
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
		_, err = w.Write([]byte(response.RawBody))
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}
