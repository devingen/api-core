package server

import (
	"github.com/devingen/api-core/dvnruntime"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func ReturnResponse(w http.ResponseWriter, response dvnruntime.Response, err error) {

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
		_, err = w.Write([]byte(response.Body))
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}

// CORSRouterDecorator applies CORS headers to a mux.Router
type CORSRouterDecorator struct {
	R *mux.Router
}

// ServeHTTP wraps the HTTP server enabling CORS headers.
// For more info about CORS, visit https://www.w3.org/TR/cors/
func (c *CORSRouterDecorator) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Type, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}

	c.R.ServeHTTP(rw, req)
}
