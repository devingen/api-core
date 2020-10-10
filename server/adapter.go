package server

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/devingen/api-core/dvnruntime"
	"github.com/devingen/api-core/util"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
)

func AdaptResponse(resp dvnruntime.Response, err error) (events.APIGatewayProxyResponse, error) {
	awsResponse := events.APIGatewayProxyResponse{
		StatusCode:      resp.StatusCode,
		Headers:         resp.Headers,
		Body:            resp.Body,
		IsBase64Encoded: resp.IsBase64Encoded,
	}
	return awsResponse, err
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

func AdaptRequest(req *http.Request) dvnruntime.Request {
	body, _ := ioutil.ReadAll(req.Body)  // TODO handle error
	ip, _ := util.GetClientIPHelper(req) // TODO handle error

	return dvnruntime.Request{
		Resource:              "not-implemented-in-server-adapter",
		Path:                  req.URL.Path,
		HTTPMethod:            req.Method,
		Headers:               convertHeaders(req.Header),
		QueryStringParameters: convertQueryParams(req.URL.Query()),
		PathParameters:        mux.Vars(req),
		StageVariables:        nil,
		RequestContext:        dvnruntime.ProxyRequestContext{},
		Body:                  string(body),
		IsBase64Encoded:       false,
		IP:                    ip,
	}
}
