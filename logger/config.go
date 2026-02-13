package logger

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	FmtJson = iota
	FmtText
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
	project   string                           //项目名称
	level     string                           //日志级别
	fmt       int                              //日志输出格式
	file      string                           //日志保存地址
	errorFile string                           //错误日志保存地址
	traceFun  func(ctx context.Context) string //返回traceId的方法
}

func initOptions(opts ...func(*options)) options {
	var option = options{
		fmt:   FmtText,
		level: LevelDebug,
	}
	for _, opt := range opts {
		opt(&option)
	}
	return option
}

var Opt struct {
	Level     func(level string) func(*options)
	Fmt       func(fmt int) func(*options)
	File      func(file string) func(*options)
	ErrorFile func(file string) func(*options)
	TraceFun  func(traceFun func(ctx context.Context) string) func(*options)
}

// InitLog
// opts入参
// traceFun 获取traceId方法  默认使用 trace.TraceIDFromContext
// Opt.Fmt(FmtText) 设置格式化类型 默认FMT_TEXT类型
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
	loggersLevel, err := logrus.ParseLevel(option.level)
	if err != nil {
		log.Fatal("logger:illegal level ", option.level)
	}
	logrus.SetLevel(loggersLevel)
	if option.fmt == FmtJson {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000", //时间格式化
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000", //时间格式化
		})
	}
	var file *os.File
	if option.file != "" {
		var err error
		file, err = os.OpenFile(option.file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Println("logger:failed to open log file", option.file, "err", err)
			logrus.SetOutput(io.MultiWriter(os.Stdout))
		} else {
			logrus.SetOutput(io.MultiWriter(os.Stdout, file))
		}
	} else {
		logrus.SetOutput(io.MultiWriter(os.Stdout))
	}

	var errorFile *os.File
	if option.errorFile != "" {
		var err error
		errorFile, err = os.OpenFile(option.errorFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Println("logger:failed to open error log file", option.errorFile, "err", err)
			errorLogger.SetOutput(io.MultiWriter(os.Stderr))
		} else {
			errorLogger.SetOutput(io.MultiWriter(os.Stderr, errorFile))
		}
	} else {
		errorLogger.SetOutput(io.MultiWriter(os.Stderr))
	}

	return func() {
		Debug(ctx, "logger:defer close logger")
		if file != nil {
			_ = file.Close()
		}
		if errorFile != nil {
			_ = errorFile.Close()
		}
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
	Opt.File = func(file string) func(*options) {
		return func(o *options) {
			o.file = file
		}
	}
	Opt.ErrorFile = func(file string) func(*options) {
		return func(o *options) {
			o.errorFile = file
		}
	}
	Opt.TraceFun = func(traceFun func(ctx context.Context) string) func(*options) {
		return func(o *options) {
			o.traceFun = traceFun
		}
	}
}
