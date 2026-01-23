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
}

type fullerror struct {
	code  int
	msg   string
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
func New(code int, msg string) Error {
	return fullerror{
		code:  code,
		msg:   msg,
		stack: callers(),
	}
}

func Wrap(code int, msg string, cause error) Error {
	return fullerror{
		code:  code,
		msg:   msg,
		stack: callers(),
		cause: cause,
	}
}

func (f fullerror) Error() string {
	s := f.TopStackTrace()
	return fmt.Sprintf("code:%d,msg:%s-%s:%d", f.code, f.msg, s, s)
}

func (f fullerror) Code() int {
	return f.code
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
