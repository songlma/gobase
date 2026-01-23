package trace

import (
	"context"
	"testing"
	"time"

	"github.com/opentracing/opentracing-go"
)

func TestJaeger(t *testing.T) {
	var ctx = context.TODO()
	closer, err := InitJaeger(Config{
		Service:  "TestJaeger",
		url:      "http://zipkin.istio-system:9411",
		hostPort: "9411",
	})
	if err != nil {
		t.Errorf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	defer closer.Close()

	span1, ctx := opentracing.StartSpanFromContext(ctx, "span_1")
	time.Sleep(time.Second / 2)

	span11, _ := opentracing.StartSpanFromContext(ctx, "span_1-1")
	time.Sleep(time.Second / 2)
	span11.Finish()

	span1.Finish()

}
