package aws

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/devingen/api-core/dto"
)

func AdaptResponse(resp dto.Response, err error) (events.APIGatewayProxyResponse, error) {
	awsResponse := events.APIGatewayProxyResponse{
		StatusCode:      resp.StatusCode,
		Headers:         resp.Headers,
		Body:            resp.Body,
		IsBase64Encoded: resp.IsBase64Encoded,
	}
	return awsResponse, err
}

//func AdaptRequest(req events.APIGatewayProxyRequest) dto.Request {
//	return dto.Request{
//		Resource:              req.Resource,
//		Path:                  req.Path,
//		HTTPMethod:            req.HTTPMethod,
//		Headers:               req.Headers,
//		QueryStringParameters: req.QueryStringParameters,
//		PathParameters:        req.PathParameters,
//		StageVariables:        req.StageVariables,
//		RequestContext: dto.ProxyRequestContext{
//			AccountID:  req.RequestContext.AccountID,
//			ResourceID: req.RequestContext.ResourceID,
//			Stage:      req.RequestContext.Stage,
//			RequestID:  req.RequestContext.RequestID,
//			Identity: dto.RequestIdentity{
//				CognitoIdentityPoolID:         req.RequestContext.Identity.CognitoIdentityPoolID,
//				AccountID:                     req.RequestContext.Identity.AccountID,
//				CognitoIdentityID:             req.RequestContext.Identity.CognitoIdentityID,
//				Caller:                        req.RequestContext.Identity.Caller,
//				APIKey:                        req.RequestContext.Identity.APIKey,
//				SourceIP:                      req.RequestContext.Identity.SourceIP,
//				CognitoAuthenticationType:     req.RequestContext.Identity.CognitoAuthenticationType,
//				CognitoAuthenticationProvider: req.RequestContext.Identity.CognitoAuthenticationProvider,
//				UserArn:                       req.RequestContext.Identity.UserArn,
//				UserAgent:                     req.RequestContext.Identity.UserAgent,
//				User:                          req.RequestContext.Identity.User,
//			},
//			ResourcePath: req.RequestContext.ResourcePath,
//			Authorizer:   req.RequestContext.Authorizer,
//			HTTPMethod:   req.RequestContext.HTTPMethod,
//			APIID:        req.RequestContext.APIID,
//		},
//		Body:            req.Body,
//		IsBase64Encoded: req.IsBase64Encoded,
//	}
//}
