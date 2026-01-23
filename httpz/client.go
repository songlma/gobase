package httpz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Client struct {
	c           *http.Client
	Opentracing bool
}

func Get(ctx context.Context, url string) (resp *http.Response, err error) {
	return NewDefaultClient().Get(ctx, url)
}
func (client *Client) Get(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(ctx, req)
}

func PostForm(ctx context.Context, url string, data url.Values, header http.Header) (resp *http.Response, err error) {
	return NewDefaultClient().PostForm(ctx, url, data, header)
}

func (client *Client) PostForm(ctx context.Context, url string, data url.Values, header http.Header) (resp *http.Response, err error) {
	return client.Post(ctx, url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()), header)
}

// PostJson
// 示例:
// resp, err := httpz.PostJson(ctx, url, params, header)
//
//	if err != nil {
//		return nil, errorz.GoErr(err)
//	}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return body, errorz.GoErr(err)
//	}
func PostJson(ctx context.Context, url string, data interface{}, header http.Header) (resp *http.Response, err error) {
	return NewDefaultClient().PostJson(ctx, url, data, header)
}

func (client *Client) PostJson(ctx context.Context, url string, data interface{}, header http.Header) (resp *http.Response, err error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return client.Post(ctx, url, "application/json", bytes.NewBuffer(jsonStr), header)
}

func PostFile(ctx context.Context, url string, readBuff io.Reader, fieldName, fileName string, header http.Header) (resp *http.Response, err error) {
	return NewDefaultClient().PostFile(ctx, url, readBuff, fieldName, fileName, header)
}

// PostFile reader bufio.NewReader
func (client *Client) PostFile(ctx context.Context, url string, readBuff io.Reader, fieldName, fileName string, header http.Header) (resp *http.Response, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fileWriter, readBuff)
	if err != nil {
		return nil, err
	}
	err = bodyWriter.Close()
	if err != nil {
		return nil, err
	}
	contentType := bodyWriter.FormDataContentType()
	return client.Post(ctx, url, contentType, bodyBuf, header)
}

func (client *Client) Post(ctx context.Context, url, contentType string, body io.Reader, header http.Header) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	for key, _ := range header {
		req.Header.Set(key, header.Get(key))
	}
	req.Header.Set("Content-Type", contentType)

	return client.Do(ctx, req)
}

func (client *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if !client.Opentracing {
		return client.c.Do(req)
	}
	operationName := fmt.Sprintf("HTTP Client %s %s", req.Method, req.URL.String())
	span, _ := opentracing.StartSpanFromContext(ctx, operationName)
	defer span.Finish()
	ext.HTTPMethod.Set(span, req.Method)
	ext.SpanKind.Set(span, ext.SpanKindRPCClientEnum)
	ext.HTTPUrl.Set(span, req.URL.String())
	ext.Component.Set(span, defaultComponentName)
	response, err := client.c.Do(req)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error")
		span.LogKV("error.kind", err.Error())
		span.LogKV("error.object", err.Error())
	} else {
		ext.HTTPStatusCode.Set(span, uint16(response.StatusCode))
		if response.StatusCode != 200 {
			ext.Error.Set(span, true)
			span.LogKV("event", "error")
			span.LogKV("error.kind", response.StatusCode)
			span.LogKV("error.object", response.StatusCode)
		}
	}
	return response, err
}

func NewClientWithHttpClient(c *http.Client) *Client {
	if c.Timeout == 0 {
		c.Timeout = 10 * time.Second
	}
	return &Client{
		c: c,
	}
}

func NewDefaultClient() *Client {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	return &Client{
		c:           client,
		Opentracing: true,
	}
}
