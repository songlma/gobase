package trace

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

//SamplerType
//probabilistic采样器以采样概率等于sampler.param属性值的方式进行随机采样决策。例如，每sampler.param=0.110 条轨迹中大约有 1 条将被采样。
//ratelimiting 采样器使用漏桶速率限制器来确保以某个恒定速率对轨迹进行采样。例如，当sampler.param=2.0它以每秒 2 次跟踪的速率对请求进行采样时。

type ZipkinConfig struct {
	Service, Url, HostPort string
	Logger                 LogInterface
}

func InitZipkin(traceCfg ZipkinConfig) (io.Closer, error) {
	reporter := zipkinhttp.NewReporter(traceCfg.Url)
	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(traceCfg.Service, traceCfg.HostPort)
	if err != nil {
		fmt.Printf("unable to create local endpoint: %+v\n", err)
	}
	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		fmt.Printf("unable to create tracer: %+v\n", err)
	}
	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)
	opentracing.SetGlobalTracer(tracer)
	if traceCfg.Logger != nil {
		logger = traceCfg.Logger
	}
	return reporter, nil
}
