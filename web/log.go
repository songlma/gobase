package web

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/songlma/gobase/contextz"
)

func errorLog(ctx context.Context, tag string, err error, args ...interface{}) {
	traceId, _ := contextz.GetTraceID(ctx)
	var argInfo []string
	for _, arg := range args {
		argInfo = append(argInfo, fmt.Sprintf("%v", arg))
	}
	log.Printf("err:%v;trace_id:%s;tag:%s;%s", err, traceId, tag, strings.Join(argInfo, ";"))
}
