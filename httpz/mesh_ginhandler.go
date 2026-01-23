package httpz

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var EnvKey = "env"
var XRequestIdKey = "x-request-id"
var XB3TraceIdKey = "x-b3-traceid"
var XB3SpanIdKey = "x-b3-spanid"
var XB3ParentSpanIdKey = "x-b3-parentspanid"
var XB3SampledKey = "x-b3-sampled"
var XB3FlagsKey = "x-b3-flags"
var XOtSpanContextKey = "x-ot-span-context"

type HTTPHeadersCarrier struct {
	XRequestId      string
	XB3TraceId      string
	XB3SpanId       string
	XB3ParentSpanId string
	XB3Sampled      string
	XB3Flags        string
	XOtSpanContext  string
	Env             string
}

var HTTPHeadersCarrierKey = "http_header_carrier"

func MeshGinHandlerFunc() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := ginCtx.Request.Context()
		var carrier = &HTTPHeadersCarrier{}
		carrier.Env = ginCtx.Request.Header.Get(EnvKey)
		carrier.XRequestId = ginCtx.Request.Header.Get(XRequestIdKey)
		carrier.XB3TraceId = ginCtx.Request.Header.Get(XB3TraceIdKey)
		carrier.XB3SpanId = ginCtx.Request.Header.Get(XB3SpanIdKey)
		carrier.XB3ParentSpanId = ginCtx.Request.Header.Get(XB3ParentSpanIdKey)
		carrier.XB3Sampled = ginCtx.Request.Header.Get(XB3SampledKey)
		carrier.XB3Flags = ginCtx.Request.Header.Get(XB3FlagsKey)
		carrier.XOtSpanContext = ginCtx.Request.Header.Get(XOtSpanContextKey)
		ctx = context.WithValue(ctx, HTTPHeadersCarrierKey, carrier)
		ginCtx.Request = ginCtx.Request.WithContext(ctx)
		ginCtx.Next()
	}
}

func AddMeshHeader(ctx context.Context, header http.Header) {
	value := ctx.Value(HTTPHeadersCarrierKey)
	if value == nil {
		fmt.Println("http_header_carrier not found")
		return
	}
	if httpCarrier, ok := value.(*HTTPHeadersCarrier); ok {
		if httpCarrier.Env != "" {
			header.Set(EnvKey, httpCarrier.Env)
		}
		if httpCarrier.XRequestId != "" {
			header.Set(XRequestIdKey, httpCarrier.XRequestId)
		}

		if httpCarrier.XB3TraceId != "" {
			header.Set(XB3TraceIdKey, httpCarrier.XB3TraceId)
		}
		if httpCarrier.XB3SpanId != "" {
			header.Set(XB3SpanIdKey, httpCarrier.XB3SpanId)
		}
		if httpCarrier.XB3ParentSpanId != "" {
			header.Set(XB3ParentSpanIdKey, httpCarrier.XB3ParentSpanId)
		}
		if httpCarrier.XB3Sampled != "" {
			header.Set(XB3SampledKey, httpCarrier.XB3Sampled)
		}
		if httpCarrier.XB3Flags != "" {
			header.Set(XB3FlagsKey, httpCarrier.XB3Flags)
		}
		if httpCarrier.XOtSpanContext != "" {
			header.Set(XOtSpanContextKey, httpCarrier.XOtSpanContext)
		}
	}
}
