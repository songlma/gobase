package web

import (
	"bytes"
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/songlma/gobase/contextz"
	"github.com/songlma/gobase/logger"
	"github.com/songlma/gobase/trace"
)

const CorralIdKey = "CorralId"

var writerPool = &sync.Pool{
	New: func() interface{} {
		return &bodyLogWriter{
			bodyBuf: bytes.NewBufferString(""),
		}
	},
}

func InitContext(gctx *gin.Context) {
	requestNanoTime := time.Now()
	requestTime := requestNanoTime.Format("2006-01-02T15:04:05.000Z07:00")
	var traceId string
	if xTraceID := gctx.GetHeader("X-Trace-ID"); len(xTraceID) > 0 {
		traceId = xTraceID
	} else {
		if xRequestID := gctx.GetHeader("X-Request-ID"); len(xRequestID) > 0 {
			traceId = xRequestID
		} else {
			traceId = strconv.FormatInt(requestNanoTime.UnixNano(), 32)
		}
	}
	ctx := gctx.Request.Context()
	ctx, _ = contextz.SetTraceID(ctx, traceId)
	ctx = context.WithValue(ctx, trace.TraceIdKey, traceId)
	ctx = context.WithValue(ctx, trace.SpanIdKey, strconv.FormatInt(requestNanoTime.UnixNano(), 10))

	corralId := gctx.Request.Header.Get(CorralIdKey)
	if corralId == "" {
		corralId = strconv.FormatInt(time.Now().UnixNano(), 32)
		gctx.Request.Header.Add(CorralIdKey, corralId)
	}
	ctx, _ = contextz.SetCorralID(ctx, corralId)
	ctx = context.WithValue(ctx, "X-Request-Time", requestTime)
	ctx = context.WithValue(ctx, "X-Request-Unixtime", requestNanoTime.UnixNano())
	gctx.Request = gctx.Request.WithContext(ctx)
	gctx.Next()
}

func RequestLog(gctx *gin.Context) {
	sTime := time.Now().UnixNano()
	bodylogWriter := writerPool.Get().(*bodyLogWriter)
	bodylogWriter.Init(gctx.Writer)
	gctx.Writer = bodylogWriter
	gctx.Next()
	params, err := GetParams(gctx)
	if err != nil {
		errorLog(gctx.Request.Context(), "GetParams", err)
	}
	ctx := gctx.Request.Context()
	corralId, _ := contextz.GetCorralID(ctx)
	path := gctx.Request.URL.Path
	if !strings.Contains(path, "inner") {
		logger.WithFields(ctx, logger.Fields{
			"type":     "income",
			"ts":       float64(time.Now().UnixNano()-sTime) / 1000000,
			"corralId": corralId,
			"method":   path,
			"params":   string(params),
			"result":   strings.Trim(bodylogWriter.BodyString(), "\n"),
		}).Infof(
			"%s|%d",
			gctx.Request.Method,
			gctx.Writer.Status(),
		)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

func (w *bodyLogWriter) Init(writer gin.ResponseWriter) {
	w.bodyBuf.Reset()
	w.ResponseWriter = writer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.bodyBuf.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) BodyString() string {
	return w.bodyBuf.String()
}
