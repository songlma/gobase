package httpz

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/songlma/gobase/trace"
)

func TestAppRequest(t *testing.T) {
	ctx := context.Background()
	closer, err := trace.InitJaeger(trace.Config{
		Service:            "Gov2",
		LocalAgentHostPort: "localhost:6831",
		LogSpans:           true,
	})
	if err != nil {
		t.Error("InitJaegerErr:", err)
	}
	defer func() {
		if closer != nil {
			closer.Close()
		}
	}()

	handler := GetGinHandler(ctx)
	params := GetRecentPlayListParams() //
	//params := GetSignErrParams() //
	req, _ := http.NewRequest("POST", "/2_8/story/recent_play_list", strings.NewReader(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(ctx)
	//构建带tracer的请求
	req = req.WithContext(ctx)
	w := newCloseNotifyingRecorder()
	handler.ServeHTTP(w, req)
	resBt, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("ioReadAllErr: %v", err)
	}
	t.Log("resBt", string(resBt))
}

func GetGinHandler(ctx context.Context) *gin.Engine {
	ginEngine := DefaultGin(nil)
	appGroup := ginEngine.Group("/2_8")
	appGroup.POST("/:client/:service", func(c *gin.Context) {

		c.JSON(http.StatusBadGateway, "ok")

	})
	return ginEngine
}

func GetRecentPlayListParams() url.Values {
	var params = url.Values{}
	params.Add("os", "android")
	params.Add("net_type", "1")
	return params
}

func GetSignErrParams() url.Values {
	var params = url.Values{}
	params.Add("user_name", "")
	params.Add("channel", "android")
	params.Add("ua", "Dalvik/2.1.0 (Linux; U; Android 10; M2007J17C MIUI/V12.0.11.0.QJSCNXM)")
	params.Add("device_uuid", "094945127aeb0c33f7107e975b0a274")
	params.Add("mac", "50:8E:49:38:3E:32")
	params.Add("screen_height", "2179")
	return params
}

type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func newCloseNotifyingRecorder() *closeNotifyingRecorder {
	return &closeNotifyingRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func (c *closeNotifyingRecorder) close() {
	c.closed <- true
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}
