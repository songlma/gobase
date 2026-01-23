package httpz

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/songlma/gobase/trace"
)

func TestName(t *testing.T) {
	closer, err := trace.InitJaeger(trace.Config{
		Service:            "Gov2",
		LocalAgentHostPort: "localhost:6831",
		LogSpans:           true,
		SamplerType:        "probabilistic",
		SamplerParam:       1,
	})
	if err != nil {
		t.Error("InitJaegerErr:", err)
	}
	defer func() {
		if closer != nil {
			closer.Close()
		}
	}()

	ctx := context.Background()
	resp, err := PostJson(ctx, "http://localhost:8080/api/test", nil, nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
}

func TestClient_PostFile(t *testing.T) {
	ctx := context.Background()
	media, err := UploadTempMedia(ctx, "")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(media)
}

func UploadTempMedia(ctx context.Context, imageUrl string) (string, error) {
	token := "54_L5SSetZjUU3l_LHwOuwdv_dhwhp3c1QhWnYfRMeJR4joMhlzTGo3QgpPWTTz28KSydlbRNcnNQ2CfaynYd2nfR8mn8L8lYxFUyzz5yEItC4BovNuO3QiUqE7abslXSQWwo-P6ipVs-htla0uNZYeAGAWIH"
	resp, err := http.Get(imageUrl)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	url := "https://api.weixin.qq.com/cgi-bin/media/upload?access_token=" + token + "&type=image"

	resp, err = PostFile(ctx, url, resp.Body, "media", "865116d126582ebeba0d799a9dc921aab.jpg", nil)

	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return string(body), nil
}
