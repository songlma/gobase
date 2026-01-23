package httpz

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

type Config struct {
	k8sReadinessSpan bool
}

func DefaultConfig() *Config {
	return &Config{
		k8sReadinessSpan: false,
	}
}
func DefaultGin(config *Config, middleware ...gin.HandlerFunc) *gin.Engine {
	ginEngine := gin.New()
	if config == nil {
		config = DefaultConfig()
	}
	var options []MWOption

	if !config.k8sReadinessSpan {
		options = append(options, IncludingSpans(func(method string) bool {
			if strings.Contains(method, "k8s_readiness") {
				return false
			} else {
				return true
			}
		}))
	}
	openTracingGinHandlerFunc := OpenTracingGinHandlerFunc(
		opentracing.GlobalTracer(), options...,
	)
	middleware = append(middleware, openTracingGinHandlerFunc, PanicGinHandlerFunc(), MeshGinHandlerFunc())
	ginEngine.Use(middleware...)
	return ginEngine
}
