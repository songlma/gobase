package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/songlma/gobase/contextz"
	"github.com/songlma/gobase/errorz"
)

var (
	loggerPackage      string
	callerInitOnce     sync.Once
	minimumCallerDepth int
)

type Fields map[string]interface{}

const (
	maximumCallerDepth int = 25
	knownLoggerFrames  int = 4
)

var errorLogger logrus.Logger

func init() {
	//默认输出
	logrus.SetOutput(io.MultiWriter(os.Stdout))
	errorLogger.SetOutput(os.Stderr)
	errorLogger.SetLevel(logrus.ErrorLevel)
	errorLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000", //时间格式化
	})

}

func Debug(ctx context.Context, args ...interface{}) {
	entry := logrus.WithContext(ctx)
	commonEntry(ctx, entry, nil).Debug(args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	entry := logrus.WithContext(ctx)
	commonEntry(ctx, entry, nil).Debugf(format, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	entry := logrus.WithContext(ctx)
	fields := map[string]interface{}{}
	var errArgs []interface{}
	for _, arg := range args {
		var stack string
		if err, ok := arg.(errorz.Error); ok {
			stack = fmt.Sprintf("%+v\n", err)
			errArgs = append(errArgs, stack)
		} else {
			errArgs = append(errArgs, arg)
		}
		if stack != "" {
			fields["stack"] = stack
			break
		}
	}
	//commonEntry(ctx, errorLogger.WithContext(ctx), nil).Error(errArgs...)
	commonEntry(ctx, entry, fields).Error(args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	entry := logrus.WithContext(ctx)
	fields := map[string]interface{}{}
	//errorLogger.Errorf(format, args...)
	for _, arg := range args {
		var stack string
		if err, ok := arg.(errorz.Error); ok {
			stack = fmt.Sprintf("%+v\n", err)
			//errorLogger.Errorf("%+v\n", err)
		}
		if stack != "" {
			fields["stack"] = stack
			break
		}
	}
	commonEntry(ctx, entry, fields).Errorf(format, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	entry := logrus.WithContext(ctx)
	commonEntry(ctx, entry, nil).Info(args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	entry := logrus.WithContext(ctx)
	commonEntry(ctx, entry, nil).Infof(format, args...)
}

func Trace(ctx context.Context, args ...interface{}) {
	entry := logrus.WithContext(ctx)
	commonEntry(ctx, entry, nil).Trace(args...)
}

func Tracef(ctx context.Context, format string, args ...interface{}) {
	entry := logrus.WithContext(ctx)
	commonEntry(ctx, entry, nil).Tracef(format, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	entry := logrus.WithContext(ctx)
	commonEntry(ctx, entry, nil).Warn(args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	entry := logrus.WithContext(ctx)
	commonEntry(ctx, entry, nil).Warnf(format, args...)
}

func WithFields(ctx context.Context, fileds map[string]interface{}) *logrus.Entry {
	entry := logrus.WithContext(ctx)
	return commonEntry(ctx, entry, fileds)
}

func commonEntry(ctx context.Context, entry *logrus.Entry, fields logrus.Fields) *logrus.Entry {
	frame := getCaller()
	traceId := ""
	if traceFun != nil {
		traceId = traceFun(ctx)
	} else if ctx != nil {
		value, err := contextz.GetTraceID(ctx)
		if err != nil {
			traceId = value
		}
	}
	commonFields := logrus.Fields{
		"svc":    svc,
		"caller": fmt.Sprintf("%v:%d", frame.Func.Name(), frame.Line),
		"type":   "all",
		"trace":  traceId,
	}
	if fields != nil {
		for k, v := range fields {
			commonFields[k] = v
		}
	}
	return entry.WithFields(commonFields)
}

// Copy from logrus
func getCaller() *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, 2)
		_ = runtime.Callers(0, pcs)
		loggerPackage = getPackageName(runtime.FuncForPC(pcs[1]).Name())

		// now that we have the cache, we can skip a minimum count of known-logrus functions
		// XXX this is dubious, the number of frames may vary
		minimumCallerDepth = knownLoggerFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		// If the caller isn't part of this package, we're done
		if pkg != loggerPackage {
			return &f
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}

func nano2Milli(nanoSecond int64) float64 {
	return float64(nanoSecond) / 1000000
}
