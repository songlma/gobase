package httpz

import (
	"bytes"
	"sync"

	"github.com/gin-gonic/gin"
)

var writerPool = &sync.Pool{
	New: func() interface{} {
		return &BodyLogWriter{
			bodyBuf: bytes.NewBufferString(""),
		}
	},
}

func GetBodyLogWriter() *BodyLogWriter {
	return writerPool.Get().(*BodyLogWriter)
}

func PutBodyLogWriter(bodyLogWriter *BodyLogWriter) {
	writerPool.Put(bodyLogWriter)
}

type BodyLogWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

func (w *BodyLogWriter) Init(writer gin.ResponseWriter) {
	w.bodyBuf.Reset()
	w.ResponseWriter = writer
}

func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.bodyBuf.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w BodyLogWriter) BodyString() string {
	return w.bodyBuf.String()
}
