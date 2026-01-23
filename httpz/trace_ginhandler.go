package httpz

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const defaultComponentName = "net/http"

type mwOptions struct {
	opNameFunc         func(r *http.Request) string
	spanCreateObserver []func(ctx context.Context, span opentracing.Span, r *http.Request)
	spanFinishObserver []func(ctx context.Context, span opentracing.Span, r *http.Request)
	urlTagFunc         func(u *url.URL) string
	componentName      string
	inclusionFunc      SpanInclusionFunc
}

// MWOption controls the behavior of the Middleware.
type MWOption func(*mwOptions)

// OperationNameFunc returns a MWOption that uses given function f
// to generate operation name for each server-side span.
func OperationNameFunc(f func(r *http.Request) string) MWOption {
	return func(options *mwOptions) {
		options.opNameFunc = f
	}
}

// MWComponentName returns a MWOption that sets the component name
// for the server-side span.
func MWComponentName(componentName string) MWOption {
	return func(options *mwOptions) {
		options.componentName = componentName
	}
}

// MWSpanObserver returns a MWOption that observe the span
// for the server-side span.
// span 创建后
func MWSpanCreateObserver(f ...func(ctx context.Context, span opentracing.Span, r *http.Request)) MWOption {
	return func(options *mwOptions) {
		options.spanCreateObserver = f
	}
}

// MWSpanObserver returns a MWOption that observe the span
// for the server-side span.
// span 调用Finish前
func MWSpanFinishObserver(f ...func(ctx context.Context, span opentracing.Span, r *http.Request)) MWOption {
	return func(options *mwOptions) {
		options.spanFinishObserver = f
	}
}

// MWURLTagFunc returns a MWOption that uses given function f
// to set the span's http.url tag. Can be used to change the default
// http.url tag, eg to redact sensitive information.
func MWURLTagFunc(f func(u *url.URL) string) MWOption {
	return func(options *mwOptions) {
		options.urlTagFunc = f
	}
}

type SpanInclusionFunc func(method string) bool

// IncludingSpans binds a IncludeSpanFunc to the options
func IncludingSpans(inclusionFunc SpanInclusionFunc) MWOption {
	return func(options *mwOptions) {
		options.inclusionFunc = inclusionFunc
	}
}

// Middleware is a gin native version of the equivalent middleware in:
//
//	https://github.com/opentracing-contrib/go-stdlib/
func OpenTracingGinHandlerFunc(tr opentracing.Tracer, options ...MWOption) gin.HandlerFunc {
	opts := mwOptions{
		opNameFunc: func(r *http.Request) string {
			return fmt.Sprintf("HTTP %s %s", r.Method, r.URL.String())
		},
		urlTagFunc: func(u *url.URL) string {
			return u.String()
		},
	}
	for _, opt := range options {
		opt(&opts)
	}
	return func(c *gin.Context) {
		if opts.inclusionFunc != nil && !opts.inclusionFunc(c.Request.URL.Path) {
			c.Next()
			return
		}
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		ctx, _ := tr.Extract(opentracing.HTTPHeaders, carrier)
		op := opts.opNameFunc(c.Request)
		span := tr.StartSpan(op, ext.RPCServerOption(ctx))
		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.SpanKind.Set(span, ext.SpanKindRPCServerEnum)
		ext.HTTPUrl.Set(span, opts.urlTagFunc(c.Request.URL))
		for _, observer := range opts.spanCreateObserver {
			observer(c.Request.Context(), span, c.Request)
		}
		// set component name, use "net/http" if caller does not specify
		componentName := opts.componentName
		if componentName == "" {
			componentName = defaultComponentName
		}
		ext.Component.Set(span, componentName)
		c.Request = c.Request.WithContext(
			opentracing.ContextWithSpan(c.Request.Context(), span))
		c.Next()
		for _, observer := range opts.spanFinishObserver {
			observer(c.Request.Context(), span, c.Request)
		}
		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
		if c.Writer.Status() != http.StatusOK {
			ext.Error.Set(span, true)
			span.LogKV("event", "error")
			span.LogKV("error.kind", c.Writer.Status())
			span.LogKV("error.object", c.Writer.Status())
		}
		span.Finish()
	}
}
