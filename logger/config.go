package logger

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	FMT_JSON = iota
	FMT_TEXT
)
const (
	LevelPanic string = "panic"
	LevelFatal string = "fatal"
	LevelError string = "error"
	LevelWarn  string = "warn"
	LevelInfo  string = "info"
	LevelDebug string = "debug"
	LevelTrace string = "trace"
)

var svc string

var traceFun func(ctx context.Context) string

// 组级别
type options struct {
	project  string                           //项目名称
	level    string                           //日志级别
	fmt      int                              //日志输出格式
	traceFun func(ctx context.Context) string //返回traceId的方法
}

func initOptions(opts ...func(*options)) options {
	var option = options{
		fmt:   FMT_TEXT,
		level: LevelDebug,
	}
	for _, opt := range opts {
		opt(&option)
	}
	return option
}

var Opt struct {
	Level    func(level string) func(*options)
	Fmt      func(fmt int) func(*options)
	TraceFun func(traceFun func(ctx context.Context) string) func(*options)
}

// InitLog
// opts入参
// traceFun 获取traceId方法  默认使用 trace.TraceIDFromContext
// Opt.Fmt(FMT_TEXT) 设置格式化类型 默认FMT_TEXT类型
// Opt.Level(LevelDebug) 设置日志级别 默认LevelDebug级别
func InitLog(ctx context.Context, project string, traceF func(ctx context.Context) string, opts ...func(*options)) func() {
	option := initOptions(opts...)
	svc = project
	option.project = project
	traceFun = traceF
	//go 标准日志
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(io.MultiWriter(os.Stderr))
	//logrus ReportCaller 务必关闭
	logrus.SetReportCaller(false)
	logrusLevel, err := logrus.ParseLevel(option.level)
	if err != nil {
		log.Fatal("logger:illegal level ", option.level)
	}
	logrus.SetLevel(logrusLevel)
	if option.fmt == FMT_JSON {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000", //时间格式化
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000", //时间格式化
		})
	}
	logrus.SetOutput(io.MultiWriter(os.Stdout))
	return func() {
		Debug(ctx, "logger:defer close logger")
	}
}

func init() {
	Opt.Level = func(level string) func(*options) {
		return func(o *options) {
			o.level = level
		}
	}
	Opt.Fmt = func(fmt int) func(*options) {
		return func(o *options) {
			o.fmt = fmt
		}
	}
	Opt.TraceFun = func(traceFun func(ctx context.Context) string) func(*options) {
		return func(o *options) {
			o.traceFun = traceFun
		}
	}
}
