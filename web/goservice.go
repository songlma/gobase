package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/songlma/gobase/trace"
)

const CodeOk int64 = 0

type Request struct {
	Version string      `json:"version"`
	Params  interface{} `json:"params"`
}

type Result struct {
	Code    int64       `json:"code"`
	Msg     string      `json:"msg"`
	Content interface{} `json:"content"`
	Alert   string      `json:"alert"`
	TraceId string      `json:"trace_id"`
}

/*
*
绑定参数
*/
func ShouldBindBodyWith(ginContext *gin.Context, obj interface{}) error {
	body, errz := GetParams(ginContext)
	if errz != nil {
		return errz
	}
	var object json.RawMessage
	var res = Request{
		Params: &object,
	}
	err := json.Unmarshal(body, &res)
	if err != nil {
		return err
	}
	if res.Version == "" {

	}
	switch res.Version {
	case "V2":
		if err = json.Unmarshal(object, obj); err != nil {
			return err
		}
	default:
		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetParams(ginCtx *gin.Context) (body []byte, err error) {
	ContentType := ginCtx.Request.Header.Get("Content-Type")
	if strings.Contains(ContentType, "application/json") {
		if cb, ok := ginCtx.Get(gin.BodyBytesKey); ok {
			if body, ok = cb.([]byte); ok {
				return body, nil
			}
		}
		body, err = ioutil.ReadAll(ginCtx.Request.Body)
		if err != nil {
			return nil, err
		}
		ginCtx.Set(gin.BodyBytesKey, body)
	} else if strings.Contains(ContentType, "application/x-www-form-urlencoded") {
		err = ginCtx.Request.ParseForm()
		if err != nil {
			errorLog(ginCtx.Request.Context(), "Request-ParseForm", err)
			return nil, err
		}
		bodyMap := map[string]string{}
		for s, i := range ginCtx.Request.Form {
			if len(i) == 0 {
				bodyMap[s] = ""
			} else {
				bodyMap[s] = i[0]
			}
		}
		body, err = json.Marshal(bodyMap)
		if err != nil {
			return nil, err
		}
	}
	return body, nil
}

func NewApiResult() *Result {
	return &Result{
		Code: -999,
	}
}

func SetApiResultSuccess(ctx context.Context, resp http.ResponseWriter, content interface{}) {
	setApiResult(ctx, resp, CodeOk, "success", "操作成功", content)
}

func SetApiResultError(ctx context.Context, resp http.ResponseWriter, err error, alert string) {
	var code = 0
	if se, ok := err.(interface {
		Code() int
	}); ok {
		code = se.Code()
	}
	setApiResult(ctx, resp, int64(code), err.Error(), alert, nil)
}

func setApiResult(ctx context.Context, resp http.ResponseWriter, resultCode int64, msg string, alert string, content interface{}) {
	result := NewApiResult()
	result.Code = resultCode
	result.Msg = msg
	result.Alert = alert
	result.Content = content
	result.TraceId = trace.TraceIDFromContext(ctx)
	showApiResult(ctx, resp, result)
}

func showApiResult(ctx context.Context, resp http.ResponseWriter, result *Result) {
	resultBytes, _ := json.Marshal(result)
	fmt.Fprint(resp, string(resultBytes))
}
