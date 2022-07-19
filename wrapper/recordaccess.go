package wrapper

import (
	"context"
	"fmt"
	core "github.com/devingen/api-core"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// RecordAccess logs the request details with Apache format.
// See http://httpd.apache.org/docs/current/mod/mod_log_config.html
func RecordAccess(f core.Controller) core.Controller {
	return func(ctx context.Context, req core.Request) (*core.Response, error) {
		clientIP := req.RemoteAddr
		if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
			clientIP = clientIP[:colon]
		}

		startTime := time.Now()

		response, err := f(ctx, req)

		finishTime := time.Now()

		record := &ApacheLogRecord{
			ip:          clientIP,
			method:      req.HTTPMethod,
			uri:         req.Resource,
			protocol:    req.Proto,
			status:      http.StatusOK,
			time:        finishTime.UTC(),
			elapsedTime: finishTime.Sub(startTime),
		}
		record.Log(os.Stderr)

		return response, err
	}
}

const (
	ApacheFormatPattern = "%s - - [%s] \"%s %d %d\" %f\n"
)

type ApacheLogRecord struct {
	ip                    string
	time                  time.Time
	method, uri, protocol string
	status                int
	responseBytes         int64
	elapsedTime           time.Duration
}

func (r *ApacheLogRecord) Log(out io.Writer) {
	timeFormatted := r.time.Format("02/Jan/2006 03:04:05")
	requestLine := fmt.Sprintf("%s %s %s", r.method, r.uri, r.protocol)
	fmt.Fprintf(out, ApacheFormatPattern, r.ip, timeFormatted, requestLine, r.status, r.responseBytes,
		r.elapsedTime.Seconds())
}
