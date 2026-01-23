package web

//type Config struct {
//	Ip   string
//	Host string
//}
//
//type soaClient struct {
//	conf Config
//}
//
//func NewSoaClient(conf Config) *soaClient {
//	return &soaClient{
//		conf: conf,
//	}
//}
//
///**
//请求内部go服务
//method 请求方法名 例如：/poster/api
//params 参数指针类型
//result 对应接口的返回值 指针类型
//	code !=0 情况需要业务端手动判断
//*/
//func (soa *soaClient) RequestV2(ctx context.Context, method string, params interface{}, result interface{}) (Result, error) {
//	var (
//		request = Request{
//			Version: "V2",
//			Params:  params,
//		}
//		object json.RawMessage
//		res    = Result{
//			Content: &object,
//		}
//		corralId interface{}
//		bytes    []byte
//	)
//	corralId = ctx.Value(CorralIdKey)
//	if corralId == nil {
//		corralId = ""
//	}
//	postBody, _ := json.Marshal(request)
//	bytes, err := httpPostJsonByte(ctx, strings.Join([]string{soa.conf.Ip + method}, ""), postBody, map[string]string{
//		"Host":      soa.conf.Host,
//		CorralIdKey: corralId.(string),
//	})
//	if err != nil {
//		errorLog(ctx, "RequestInteriorGoServiceV2", err)
//		return res, err
//	}
//
//	err = json.Unmarshal(bytes, &res)
//	if err != nil {
//		errorLog(ctx, "RequestInteriorGoServiceV2-Unmarshal", err)
//		return res, err
//	}
//	if res.Code != 0 {
//		return res, nil
//	}
//	if err = json.Unmarshal(object, result); err != nil {
//		errorLog(ctx, "RequestInteriorGoServiceV2-Unmarshal", err)
//		return res, err
//	}
//	return res, nil
//}
//
///**
//获取内部请求Request
//*/
//func NewInnerPostJsonRequestV2(ctx context.Context, url, host string, params interface{}) (*http.Request, error) {
//	var (
//		request = Request{
//			Version: "V2",
//			Params:  params,
//		}
//		corralId interface{}
//	)
//	corralId = ctx.Value(CorralIdKey)
//	if corralId == nil {
//		corralId = ""
//	}
//	postBody, _ := json.Marshal(request)
//
//	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(postBody))
//	if err != nil {
//		return nil, err
//	}
//	if host != "" {
//		req.Host = host
//	}
//	req.Header.Set(CorralIdKey, corralId.(string))
//	req.Header.Set("Content-Type", "application/json")
//	return req, nil
//}
//func Do(ctx context.Context, req *http.Request, result interface{}) (Result, error) {
//
//	var (
//		object json.RawMessage
//		res    = Result{
//			Content: &object,
//		}
//		body []byte
//	)
//	client := &http.Client{
//		Transport: &http.Transport{
//			DialContext: func(ctx context.Context, netw, addr string) (net.Conn, error) {
//				conn, err := net.DialTimeout(netw, addr, 10*time.Second) //设置建立连接超时
//				if err != nil {
//					return nil, err
//				}
//				conn.SetDeadline(time.Now().Add(10 * time.Second))
//				return conn, nil
//			},
//		},
//	}
//	resp, err := client.Do(req)
//	if err != nil {
//		return res, err
//	}
//	defer func() {
//		if resp != nil {
//			resp.Body.Close()
//		}
//	}()
//
//	if resp.StatusCode != http.StatusOK {
//		return res, errors.New(resp.Status)
//	}
//
//	body, err = ioutil.ReadAll(resp.Body)
//	if err != nil {
//		errorLog(ctx, "RequestInteriorGoServiceV2", err)
//		return res, err
//	}
//	err = json.Unmarshal(body, &res)
//	if err != nil {
//		errorLog(ctx, "RequestInteriorGoServiceV2-Unmarshal", err)
//		return res, err
//	}
//	if res.Code != 0 {
//		return res, nil
//	}
//	if err = json.Unmarshal(object, result); err != nil {
//		errorLog(ctx, "RequestInteriorGoServiceV2-Unmarshal", err)
//		return res, err
//	}
//	return res, nil
//}
