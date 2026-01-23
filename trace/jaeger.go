package trace

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/uber/jaeger-client-go"
	uber "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var logger LogInterface = DefaultLogger{}

//SamplerType
//probabilistic采样器以采样概率等于sampler.param属性值的方式进行随机采样决策。例如，每sampler.param=0.110 条轨迹中大约有 1 条将被采样。
//ratelimiting 采样器使用漏桶速率限制器来确保以某个恒定速率对轨迹进行采样。例如，当sampler.param=2.0它以每秒 2 次跟踪的速率对请求进行采样时。

type Config struct {
	Service, LocalAgentHostPort string
	Logger                      LogInterface
	LogSpans                    bool
	SamplerType                 string
	SamplerParam                float64
}

func InitJaeger(traceCfg Config) (io.Closer, error) {
	if traceCfg.SamplerType == "" {
		traceCfg.SamplerType = "probabilistic"
	}
	if traceCfg.SamplerParam == 0 {
		traceCfg.SamplerParam = 0.1
	}
	cfg := &config.Configuration{
		ServiceName: traceCfg.Service,
		Sampler: &config.SamplerConfig{
			Type:  traceCfg.SamplerType,
			Param: traceCfg.SamplerParam,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           traceCfg.LogSpans,
			LocalAgentHostPort: traceCfg.LocalAgentHostPort,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	if traceCfg.Logger != nil {
		logger = traceCfg.Logger
	}
	return closer, nil
}

// TraceIDFromContext 获取 TraceID
func TraceIDFromContext(ctx context.Context) string {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		val := ctx.Value(TraceIdKey)
		if val == nil {
			return ""
		} else {
			return val.(string)
		}
	}
	if u, ok := span.Context().(uber.SpanContext); ok {
		return u.TraceID().String()
	}
	if u, ok := span.Context().(zipkinot.SpanContext); ok {
		return u.TraceID.String()
	}
	return ""
}

// 设置traceId
func ContextWithTrace(ctx context.Context, trace string) context.Context {
	return context.WithValue(ctx, TraceIdKey, trace)
}
