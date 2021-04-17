package log

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
)

// loggerKey type is an opaque type for the logger lookup in a given context.
type loggerKey struct{}

// ErrLoggerNotFound means that the logger is not in that context.
var ErrLoggerNotFound = errors.New("logger not found in context")

// WithLogger function creates a new context from a given context and associates the current version of the logger object.
func WithLogger(ctx context.Context, logger *logrus.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// Of function extracts a valid Logger object from a given context. Returns an error otherwise.
func Of(ctx context.Context) (*logrus.Logger, error) {
	logger, ok := ctx.Value(loggerKey{}).(*logrus.Logger)
	if !ok {
		return nil, ErrLoggerNotFound
	}

	return logger, nil
}
