package wrapper

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	core "github.com/devingen/api-core"
	"strings"
)

type AWSLambdaHandler = func(ctx context.Context, awsReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// WithLambdaHandler wraps the controller with the AWS Lambda func.
func WithLambdaHandler(ctx context.Context, f core.Controller) AWSLambdaHandler {

	return func(ctx context.Context, awsReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// convert AWS Lambda request to our custom request
		req := adaptAWSLambdaRequest(awsReq)

		// execute function
		response, err := f(ctx, req)

		// convert response to our custom response
		httpResponse, err := buildHTTPResponse(response, err)

		// write response data
		return adaptAWSLambdaResponse(httpResponse, err)
	}
}

func adaptAWSLambdaRequest(req events.APIGatewayProxyRequest) core.Request {
	return core.Request{
		Resource:              req.Resource,
		Path:                  req.Path,
		HTTPMethod:            req.HTTPMethod,
		Headers:               convertHeaderNamesToLowercase(req.Headers),
		QueryStringParameters: req.MultiValueQueryStringParameters,
		PathParameters:        req.PathParameters,
		StageVariables:        req.StageVariables,
		RequestContext: core.ProxyRequestContext{
			AccountID:  req.RequestContext.AccountID,
			ResourceID: req.RequestContext.ResourceID,
			Stage:      req.RequestContext.Stage,
			RequestID:  req.RequestContext.RequestID,
			Identity: core.RequestIdentity{
				CognitoIdentityPoolID:         req.RequestContext.Identity.CognitoIdentityPoolID,
				AccountID:                     req.RequestContext.Identity.AccountID,
				CognitoIdentityID:             req.RequestContext.Identity.CognitoIdentityID,
				Caller:                        req.RequestContext.Identity.Caller,
				APIKey:                        req.RequestContext.Identity.APIKey,
				SourceIP:                      req.RequestContext.Identity.SourceIP,
				CognitoAuthenticationType:     req.RequestContext.Identity.CognitoAuthenticationType,
				CognitoAuthenticationProvider: req.RequestContext.Identity.CognitoAuthenticationProvider,
				UserArn:                       req.RequestContext.Identity.UserArn,
				UserAgent:                     req.RequestContext.Identity.UserAgent,
				User:                          req.RequestContext.Identity.User,
			},
			ResourcePath: req.RequestContext.ResourcePath,
			Authorizer:   req.RequestContext.Authorizer,
			HTTPMethod:   req.RequestContext.HTTPMethod,
			APIID:        req.RequestContext.APIID,
		},
		Body:            req.Body,
		IsBase64Encoded: req.IsBase64Encoded,
		IP:              req.RequestContext.Identity.SourceIP,
	}
}

func convertHeaderNamesToLowercase(header map[string]string) map[string]string {
	headers := map[string]string{}
	for k, v := range header {
		headers[strings.ToLower(k)] = v
	}
	return headers
}

func adaptAWSLambdaResponse(resp core.Response, err error) (events.APIGatewayProxyResponse, error) {
	awsResponse := events.APIGatewayProxyResponse{
		StatusCode:      resp.StatusCode,
		Headers:         resp.Headers,
		Body:            resp.RawBody,
		IsBase64Encoded: resp.IsBase64Encoded,
	}
	return awsResponse, err
}
