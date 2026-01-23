package redisz

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/songlma/gobase/contextz"
)

func errorLog(ctx context.Context, tag string, err error, args ...interface{}) {
	traceId, _ := contextz.GetTraceID(ctx)
	var argInfo []string
	for _, arg := range args {
		argInfo = append(argInfo, fmt.Sprintf("%v", arg))
	}
	log.Printf("err:%v;trace_id:%s;tag:%s;%s", err, traceId, tag, strings.Join(argInfo, ";"))
	span, _ := opentracing.StartSpanFromContext(ctx, tag)
	ext.Error.Set(span, true)
	ext.Component.Set(span, "redis")
	span.LogKV("event", "error")
	span.LogKV("error.kind", "redis")
	span.LogKV("error.object", err.Error())
	span.LogKV("message", fmt.Sprintf("%v", err))
	span.Finish()
}
