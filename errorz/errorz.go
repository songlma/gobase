package errorz

import (
	"errors"
	"fmt"
	"io"

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

type fullerror struct {
	code  int
	msg   string
	alert string
	cause error
	*stack
}

func GoErr(cause error) Error {
	return fullerror{
		code:  -1,
		stack: callers(),
		cause: cause,
	}
}
func New(code int, msg string, alert ...string) Error {
	var a string
	if len(alert) > 0 {
		a = alert[0]
	}
	return fullerror{
		code:  code,
		msg:   msg,
		alert: a,
		stack: callers(),
	}
}

func Wrap(code int, msg string, cause error, alert ...string) Error {
	var a string
	if len(alert) > 0 {
		a = alert[0]
	}
	return fullerror{
		code:  code,
		msg:   msg,
		alert: a,
		stack: callers(),
		cause: cause,
	}
}

func WrapWithAlert(err Error, alert string) Error {
	var f fullerror
	if errors.As(err, &f) {
		f.alert = alert
		return f
	}
	return fullerror{
		code:  err.Code(),
		msg:   err.Error(),
		alert: alert,
		cause: err,
		stack: callers(),
	}
}

func NewAlertError(code int, msg string, alert string) Error {
	return fullerror{
		code:  code,
		msg:   msg,
		alert: alert,
		stack: callers(),
	}
}

func (f fullerror) Error() string {
	s := f.TopStackTrace()
	return fmt.Sprintf("code:%d,msg:%s,alert:%s-%s:%d", f.code, f.msg, f.alert, s, s)
}

func (f fullerror) Code() int {
	return f.code
}

func (f fullerror) Alert() string {
	return f.alert
}

func (f fullerror) Cause() error {
	return f.cause
}

// 支持grpc转换
func (f fullerror) GRPCStatus() *status.Status {
	return status.New(codes.Code(f.code), f.msg)
}

// go 1.13 Unwrapping
func (f fullerror) Unwrap() error {
	return f.cause
}

func (f fullerror) Is(err error) bool {
	var errz Error
	if errors.As(err, &errz) {
		return f.Code() == errz.Code()
	}
	return false
}

func (f fullerror) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, fmt.Sprintf("{code=%d,msg=%s", f.code, f.msg))
			if f.cause != nil {
				io.WriteString(s, fmt.Sprintf(",wrap=%v", f.cause))
			}
			io.WriteString(s, "}")
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, fmt.Sprintf("{code=%d,msg=%s", f.code, f.msg))
		if f.cause != nil {
			io.WriteString(s, fmt.Sprintf(",wrap=%v", f.cause))
		}
		io.WriteString(s, "}")
	case 'q':
		fmt.Fprintf(s, "%q", fmt.Sprintf("%v", f.cause))
	}
}
