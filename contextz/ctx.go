package contextz

import (
	"context"
	"errors"
)

var ValueNotSet = errors.New("context value not set")
var ContextIsNil = errors.New("context value not set")

func SetTraceID(ctx context.Context, traceID string) (context.Context, error) {
	if ctx == nil {
		return ctx, ContextIsNil
	}
	return context.WithValue(ctx, "traceid", traceID), nil
}

func GetTraceID(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", ContextIsNil
	}
	val := ctx.Value("traceid")
	if val == nil {
		return "", ValueNotSet
	} else {
		return val.(string), nil
	}
}

func SetUID(ctx context.Context, uid string) (context.Context, error) {
	if ctx == nil {
		return ctx, ContextIsNil
	}
	return context.WithValue(ctx, "uid", uid), nil
}

func GetUID(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", ContextIsNil
	}
	val := ctx.Value("uid")
	if val == nil {
		return "", ValueNotSet
	} else {
		return val.(string), nil
	}
}

func SetCorralID(ctx context.Context, corralId string) (context.Context, error) {
	if ctx == nil {
		return ctx, ContextIsNil
	}
	return context.WithValue(ctx, "CorralId", corralId), nil
}

func GetCorralID(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", ContextIsNil
	}
	val := ctx.Value("CorralId")
	if val == nil {
		return "", ValueNotSet
	} else {
		return val.(string), nil
	}
}
