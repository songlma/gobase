package errorz

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Error interface {
	Error() string
	// Unwrap go 1.13 Unwrapping
	Unwrap() error
	Cause() error
	Code() int
	Alert() string
}

type fuller struct {
	code  int
	msg   string
	alert string
	cause error
	stack []uintptr
}

// New 创建基础业务错误（自动捕获堆栈）
func New(code int, msg string, opts ...Option) Error {
	e := &fuller{
		code:  code,
		msg:   msg,
		stack: captureStack(3),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Wrap 包装现有错误，添加上下文（保留原始错误链）
func Wrap(err error, code int, msg string, opts ...Option) Error {
	if err == nil {
		return nil
	}
	e := &fuller{
		code:  code,
		msg:   msg,
		cause: err,
		stack: captureStack(3),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// FromStd 将标准 error 转为 errorz.Error（保留原始信息）
func FromStd(err error, opts ...Option) Error {
	if err == nil {
		return nil
	}
	// 若已是 errorz.Error，直接返回（避免重复包装）
	var e Error
	if errors.As(err, &e) {
		// 应用新选项（如替换 alert）
		if len(opts) > 0 {
			e = Wrap(err, e.Code(), e.Error(), opts...)
		}
		return e
	}
	f := &fuller{
		code:  -1, // 未知错误码
		msg:   err.Error(),
		cause: err,
		stack: captureStack(3),
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func captureStack(skip int) []uintptr {
	const depth = 20
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	return pcs[:n]
}

// ---------- 选项模式（安全修改不可变错误） ----------

type Option func(*fuller)

// WithAlert 设置用户提示（创建新错误时使用）
func WithAlert(alert string) Option {
	return func(e *fuller) { e.alert = alert }
}

// WithCode 覆盖错误码（谨慎使用）
func WithCode(code int) Option {
	return func(e *fuller) { e.code = code }
}

func (f fuller) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		// 1. 输出当前错误信息（含文件位置）
		_, _ = io.WriteString(s, f.Error())
		// 2. 递归输出 Cause
		if f.cause != nil {
			_, _ = fmt.Fprintf(s, "-->%v", f.cause)
		}
		// 3. 仅在 %+v 时输出堆栈
		if s.Flag('+') && f.stack != nil {
			_, _ = io.WriteString(s, "(")
			for i, pc := range f.stack {
				if i > 0 {
					_, _ = io.WriteString(s, " -> ")
				}
				fn := runtime.FuncForPC(pc)
				if fn == nil {
					_, _ = io.WriteString(s, "unknown")
					continue
				}
				file, line := fn.FileLine(pc)
				name := fn.Name()

				// 简化函数名
				nameParts := strings.Split(name, "/")
				if len(nameParts) > 2 {
					name = strings.Join(nameParts[len(nameParts)-2:], "/")
				}

				// 简化文件名
				fileParts := strings.Split(file, "/")
				if len(fileParts) > 2 {
					file = strings.Join(fileParts[len(fileParts)-2:], "/")
				}

				_, _ = fmt.Fprintf(s, " %s:%s:%d ", name, file, line)
			}
			_, _ = io.WriteString(s, ")")
		}
	case 's', 'q':
		_, _ = io.WriteString(s, f.Error())
	}
}

func (f fuller) Error() string {
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "[code=%d,msg=%s", f.code, f.msg)
	if f.alert != "" {
		_, _ = fmt.Fprintf(&sb, ",alert=%s", f.alert)
	}
	// 添加文件位置信息
	if len(f.stack) > 0 {
		pc := f.stack[0]
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			file, line := fn.FileLine(pc)
			// 简化文件名：保留最后两级
			parts := strings.Split(file, "/")
			if len(parts) > 1 {
				file = strings.Join(parts[len(parts)-1:], "/")
			}
			_, _ = fmt.Fprintf(&sb, ",ln=%s:%d", file, line)
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func (f fuller) Code() int {
	return f.code
}

func (f fuller) Alert() string {
	return f.alert
}

func (f fuller) Cause() error {
	return f.cause
}

// GRPCStatus 支持grpc转换
func (f fuller) GRPCStatus() *status.Status {
	return status.New(codes.Code(f.code), f.msg)
}

// Unwrap go 1.13 Unwrapping
func (f fuller) Unwrap() error {
	return f.cause
}

// Cause 返回错误链根因（最内层错误）
func Cause(err error) error {
	for {
		ue := errors.Unwrap(err)
		if ue == nil {
			return err
		}
		err = ue
	}
}

// CodeOf 返回错误链中首个有效业务码（跳过 -1）
func CodeOf(err error) int {
	for {
		var e Error
		if errors.As(err, &e) && e.Code() != -1 {
			return e.Code()
		}
		if ue := errors.Unwrap(err); ue != nil {
			err = ue
		} else {
			return -1
		}
	}
}

func AlertOf(err error) string {
	for {
		var e Error
		if errors.As(err, &e) && e.Alert() != "" {
			return e.Alert()
		}
		if ue := errors.Unwrap(err); ue != nil {
			err = ue
		} else {
			return ""
		}
	}
}

func (f fuller) Is(err error) bool {
	var errz Error
	if errors.As(err, &errz) {
		return f.Code() == errz.Code()
	}
	return false
}
