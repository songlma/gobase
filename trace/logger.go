package trace

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Interface logger interface
type LogInterface interface {
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	WithFields(ctx context.Context, fields map[string]interface{}) *logrus.Entry
}

type DefaultLogger struct {
}

func (l DefaultLogger) Info(context.Context, string, ...interface{}) {

}

func (l DefaultLogger) Warn(context.Context, string, ...interface{}) {

}

func (l DefaultLogger) Error(context.Context, string, ...interface{}) {

}

func (l DefaultLogger) WithFields(ctx context.Context, fields map[string]interface{}) *logrus.Entry {
	return logrus.WithContext(ctx).WithFields(fields)
}
