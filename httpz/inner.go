package httpz

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetInnerRequestParams(ginCtx *gin.Context) (body []byte, err error) {
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
	}
	return body, nil
}

type Request struct {
	Version string      `json:"version"`
	Params  interface{} `json:"params"`
}

/*
*
绑定参数
*/
func ShouldBindBodyWith(ginContext *gin.Context, obj interface{}) error {
	body, errz := GetInnerRequestParams(ginContext)
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
	//兼容老版本请求
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
