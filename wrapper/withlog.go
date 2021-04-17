package wrapper

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/log"
	"github.com/devingen/api-core/util"
	"github.com/sirupsen/logrus"
)

// WithLogger wraps the controller func by adding the logger to the context.
func WithLogger(level string, f core.Controller) core.Controller {
	return func(ctx context.Context, req core.Request) (interface{}, int, error) {

		// create logger
		logger := logrus.New().WithFields(logrus.Fields{
			"request-id": util.GenerateUUID(),
		}).Logger

		logrusLevel, parseLevelErr := logrus.ParseLevel(level)
		if parseLevelErr == nil {
			logger.SetLevel(logrusLevel)
		}

		// add logger to the context
		ctxWithLogger := log.WithLogger(ctx, logger)

		// execute function
		result, status, err := f(ctxWithLogger, req)

		if err != nil {
			logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("failed-in-controller")
			// don't return, let it build a proper error response
		}

		return result, status, err
	}
}
