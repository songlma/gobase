## 项目集成

配置文件

```yaml
trace:
  host_hort: :6831
  sampler_param: 1
  sampler_type: probabilistic
```

初始化

```go
//初始化trace追踪
    var closer io.Closer
    agentHostPort := config.GetString("config.trace.host_hort")
    samplerParam := config.GetFloat64("config.trace.sampler_param")
    samplerType := config.GetString("config.trace.sampler_type")
    if agentHostPort != "" {
        closer, err = trace.InitJaeger(trace.Config{
            Service:            serviceName,
            LocalAgentHostPort: agentHostPort,
            LogSpans:           true,
            Logger:             CommonLogger{},
            SamplerParam:       samplerParam,
            SamplerType:        samplerType,
        })
        if err != nil {
            logger.Error(ctx, "InitJaegerErr:", err)
        }
    } else {
        logger.Warn(ctx, "agentHostPort empty")
    }
    defer func() {
        if closer != nil {
            closer.Close()
        }
    }()
```

//提供日志组件

```go
type CommonLogger struct {
}
 
func (l CommonLogger) Info(ctx context.Context, info string, args ...interface{}) {
    logger.Info(ctx, info, args)
}
 
func (l CommonLogger) Warn(ctx context.Context, info string, args ...interface{}) {
    logger.Warn(ctx, info, args)
}
 
func (l CommonLogger) Error(ctx context.Context, info string, args ...interface{}) {
    logger.Error(ctx, info, args)
}
 
func (l CommonLogger) WithFields(ctx context.Context, fields map[string]interface{}) *logrus.Entry {
    return logger.WithFields(ctx, fields)
}
```

## 本地测试

https://github.com/jaegertracing/jaeger/releases/tag/v1.22.0

下载对应系统的包

执行

```shell
./jaeger-all-in-one
```

访问
http://127.0.0.1:16686/search
查看结果

初始化

```go
closer, err := InitJaeger("TestJaeger", "127.0.0.1:6831")
    if err != nil {
        t.Errorf("Could not initialize jaeger tracer: %s", err.Error())
        return
    }
    defer closer.Close()
```

获取全局Tracer

```go
opentracing.GlobalTracer()
```

创建新的span 一定要调用Finish方法，分特殊情况创建完使用defer 调用Finish

```go
span, ctx := opentracing.StartSpanFromContext(ctx, "check_user_info")
//函数结束后调用Finish方法
defer span.Finish()
```

参考文档：

[opentracing-go](https://github.com/opentracing/opentracing-go)  
[OpenTracing语义标准](https://github.com/opentracing-contrib/opentracing-specification-zh/blob/master/specification.md)