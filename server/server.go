package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	CORSAccessControlAllowHeaders             = "Access-Control-Allow-Headers"
	CORSAccessControlAllowHeadersDefaultValue = "Accept, Accept-Language, Content-Type, Authorization"
	CORSAccessControlAllowMethods             = "Access-Control-Allow-Methods"
	CORSAccessControlAllowMethodsDefaultValue = "POST, GET, OPTIONS, PUT, DELETE"
)

// CORSRouterDecorator applies CORS headers to a mux.Router
type CORSRouterDecorator struct {
	R                 *mux.Router
	Headers           map[string]string
	AllowSenderOrigin bool
}

// ServeHTTP wraps the HTTP server enabling CORS headers.
// For more info about CORS, visit https://www.w3.org/TR/cors/
func (c *CORSRouterDecorator) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if c.AllowSenderOrigin {
		if origin := req.Header.Get("Origin"); origin != "" {
			rw.Header().Set("Access-Control-Allow-Origin", origin)
		}
	}

	if c.Headers != nil {
		for key, value := range c.Headers {
			rw.Header().Set(key, value)
		}
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}

	c.R.ServeHTTP(rw, req)
}
