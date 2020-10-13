package dvnruntime

import "strings"

// RequestIdentity contains identity information for the request caller.
type RequestIdentity struct {
	CognitoIdentityPoolID         string `json:"cognitoIdentityPoolId"`
	AccountID                     string `json:"accountId"`
	CognitoIdentityID             string `json:"cognitoIdentityId"`
	Caller                        string `json:"caller"`
	APIKey                        string `json:"apiKey"`
	SourceIP                      string `json:"sourceIp"`
	CognitoAuthenticationType     string `json:"cognitoAuthenticationType"`
	CognitoAuthenticationProvider string `json:"cognitoAuthenticationProvider"`
	UserArn                       string `json:"userArn"`
	UserAgent                     string `json:"userAgent"`
	User                          string `json:"user"`
}

// ProxyRequestContext contains the information to identify the AWS account and resources invoking the
// Lambda function. It also includes Cognito identity information for the caller.
type ProxyRequestContext struct {
	AccountID    string                 `json:"accountId"`
	ResourceID   string                 `json:"resourceId"`
	Stage        string                 `json:"stage"`
	RequestID    string                 `json:"requestId"`
	Identity     RequestIdentity        `json:"identity"`
	ResourcePath string                 `json:"resourcePath"`
	Authorizer   map[string]interface{} `json:"authorizer"`
	HTTPMethod   string                 `json:"httpMethod"`
	APIID        string                 `json:"apiId"` // The API Gateway rest API Id
}

// ProxyRequest contains data coming from the API Gateway proxy
type Request struct {
	Resource              string              `json:"resource"` // The resource path defined in API Gateway
	Path                  string              `json:"path"`     // The url path for the caller
	HTTPMethod            string              `json:"httpMethod"`
	Headers               map[string]string   `json:"headers"`
	QueryStringParameters map[string]string   `json:"queryStringParameters"`
	PathParameters        map[string]string   `json:"pathParameters"`
	StageVariables        map[string]string   `json:"stageVariables"`
	RequestContext        ProxyRequestContext `json:"requestContext"`
	Body                  string              `json:"body"`
	IsBase64Encoded       bool                `json:"isBase64Encoded,omitempty"`
	IP                    string              `json:"ip"`
}

func (r *Request) GetHeader(key string) (string, bool) {
	value, hasKey := r.Headers[key]
	if !hasKey {
		value, hasKey = r.Headers[strings.ToLower(key)]
	}
	return value, hasKey
}
