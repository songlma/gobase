package httpz

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/songlma/gobase/logger"
)

func PanicGinHandlerFunc() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				ctx := ginCtx.Request.Context()
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				logger.WithFields(ctx, logger.Fields{"type": "panic"}).Errorf("GRPC: panic running job: %v\n%s", r, buf)
				span, _ := opentracing.StartSpanFromContext(ctx, "ginPanic")
				ext.Error.Set(span, true)
				span.LogKV("event", "error")
				span.LogKV("error.kind", "ginPanic")
				span.LogKV("error.object", "Panic")
				span.LogKV("message", fmt.Sprintf("%v", r))
				span.LogKV("stack", string(buf))
				span.Finish()
			}
		}()
		ginCtx.Next()
	}
}
