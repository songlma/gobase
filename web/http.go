package web

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func httpPostJsonByte(ctx context.Context, url string, params []byte, header map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}
	for key, value := range header {
		req.Header.Set(key, value)
	}
	if host, ok := header["Host"]; ok {
		req.Host = host
	}
	req.Header.Set("Content-Type", "application/json")

	return httpRequest(ctx, req)
}

func httpRequest(ctx context.Context, req *http.Request) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, 10*time.Second) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(10 * time.Second))
				return conn, nil
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
