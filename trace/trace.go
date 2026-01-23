package trace

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const TraceIdKey = "TraceId"
const SpanIdKey = "SpanId"
const ParentSpanIDKey = "ParentSpanId"

/**
logger添加TraceId信息
*/
//func GrpcServiceLoggerExtraFieldFun() func(ctx context.Context) map[string]interface{} {
//	return func(ctx context.Context) map[string]interface{} {
//		if ctx == nil {
//			return nil
//		}
//		m := map[string]interface{}{}
//		spanId := ctx.Value(SpanIdKey)
//		if spanId != nil {
//			m[SpanIdKey] = spanId
//		}
//		ParentSpanId := ctx.Value(ParentSpanIDKey)
//		if ParentSpanId != nil {
//			m[ParentSpanIDKey] = ParentSpanId
//		}
//		TraceId := ctx.Value(TraceIdKey)
//		if TraceId != nil {
//			m[TraceIdKey] = TraceId
//		}
//		return m
//	}
//}

/*
*
grpc请求添加TraceId ParentSpanId信息
*/
func GrpcOutgoingHeader() func(context.Context, metadata.MD) {
	return func(ctx context.Context, header metadata.MD) {
		//TraceId
		traceId := ctx.Value(TraceIdKey)
		if traceId != nil {
			header.Set(TraceIdKey, traceId.(string))
		}
		//spanID
		spanID := ctx.Value(SpanIdKey)
		if spanID != nil {
			header.Set(ParentSpanIDKey, spanID.(string))
		}
	}
}

/*
*
GRPC Server接收到到请求时

	将header的trace添加到context
	创建SpanId
*/
func GrpcServiceIncomingHeaderInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			traceId := md.Get(TraceIdKey)
			if len(traceId) > 0 {
				ctx = context.WithValue(ctx, TraceIdKey, traceId[0])
			}
			parentSpanId := md.Get(ParentSpanIDKey)
			if len(parentSpanId) > 0 {
				ctx = context.WithValue(ctx, ParentSpanIDKey, parentSpanId[0])
			}
		} else {
			logger.Warn(ctx, "GrpcServiceIncomingHeaderInterceptor metadata not found")
		}
		ctx = context.WithValue(ctx, SpanIdKey, strconv.FormatInt(time.Now().UnixNano(), 10))
		return handler(ctx, req)
	}
}
