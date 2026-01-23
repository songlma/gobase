package httpz

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/songlma/gobase/logger"
	"github.com/songlma/gobase/web"
)

const sign = "sign"
const timestamp = "timestamp"
const serviceName = "Service-Name"
const traceId = "Trace-ID"

/*
*
内部请求
*/
func InterSignGinHandlerFunc() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ContentType := ginCtx.Request.Header.Get("Content-Type")
		if !strings.Contains(ContentType, "application/json") {
			ginCtx.JSON(http.StatusBadRequest, "bad request Content-Type err")
			ginCtx.Abort()
			return
		}
		signHeader := ginCtx.Request.Header.Get(sign)
		timestampHeader := ginCtx.Request.Header.Get(timestamp)
		if signHeader == "" || timestampHeader == "" || timestampHeader != signHeader {
			ginCtx.JSON(http.StatusBadRequest, "bad request sign")
			ginCtx.Abort()
			return
		}
		ginCtx.Next()
	}
}

/*
*
app请求日志
*/
func InterRequestLogGinHandlerFunc() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		sTime := time.Now().UnixNano()
		logWriter := GetBodyLogWriter()
		logWriter.Init(ginCtx.Writer)
		ginCtx.Writer = logWriter
		serviceNameHeader := ginCtx.Request.Header.Get(serviceName)
		traceIdHeader := ginCtx.Request.Header.Get(traceId)
		//兼容老版本
		if traceIdHeader == "" {
			traceIdHeader = ginCtx.Request.Header.Get(web.CorralIdKey)
		}
		ginCtx.Next()
		ctx := ginCtx.Request.Context()
		path := ginCtx.Request.URL.Path
		params, _ := GetInnerRequestParams(ginCtx)
		//兼容老版本
		fields := logger.Fields{
			"type":         "inner",
			"ts":           float64(time.Now().UnixNano()-sTime) / 1000000,
			"path":         path,
			"params":       string(params),
			"Method":       ginCtx.Request.Method,
			"service_name": serviceNameHeader,
			"status":       ginCtx.Writer.Status(),
			traceId:        traceIdHeader,
		}
		logger.WithFields(ctx, fields).Info(strings.Trim(logWriter.BodyString(), "\n"))
		PutBodyLogWriter(logWriter)
	}
}

// InnerRequestSpanFinishObserver 初始span创建时 传入的值
func InnerRequestSpanFinishObserver() func(ctx context.Context, span opentracing.Span, r *http.Request) {
	return func(ctx context.Context, span opentracing.Span, r *http.Request) {
		serviceNameHeader := r.Header.Get(serviceName)
		span.SetTag(serviceName, serviceNameHeader)
	}
}
