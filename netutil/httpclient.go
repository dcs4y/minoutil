package netutil

import (
	"bytes"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type httpMethod string

const (
	HttpMethodGet  httpMethod = http.MethodGet
	HttpMethodPost httpMethod = http.MethodPost
)

type httpContentType string

const (
	HttpContentTypeEncoded httpContentType = "application/x-www-form-urlencoded"
	HttpContentTypeForm    httpContentType = "multipart/form-data"
	HttpContentTypeJson    httpContentType = "application/json"
	HttpContentTypeXml     httpContentType = "application/xml"
)

// HttpClient http请求客户端
type HttpClient struct {
	url       string              //请求地址
	transport *http.Transport     //https证书
	header    map[string][]string //header
	param     map[string][]string //表单参数
	body      io.Reader           //实际请求数据。设置后param失效。
	timeout   int64               //超时时间(秒)
	error     error               //构建过程中的错误信息
}

// HttpResponse http返回结果
type HttpResponse struct {
	StatusCode    int                 //返回状态码。e.g. 200
	Header        map[string][]string //header
	Body          []byte              //返回内容
	ContentLength int64               //返回内容长度
}

func NewHttpClient(url string) *HttpClient {
	return &HttpClient{
		url:     url,
		timeout: 15,
	}
}

func (hc *HttpClient) SetTransport(transport *http.Transport) *HttpClient {
	hc.transport = transport
	return hc
}

func (hc *HttpClient) SetHeaders(header map[string][]string) *HttpClient {
	hc.header = header
	return hc
}

func (hc *HttpClient) SetHeader(key, value string) *HttpClient {
	if hc.header == nil {
		hc.header = map[string][]string{}
	}
	if "" != value {
		hc.header[key] = []string{value}
	} else {
		delete(hc.header, key)
	}
	return hc
}

func (hc *HttpClient) AddHeader(key, value string) *HttpClient {
	if hc.header == nil {
		hc.header = map[string][]string{}
	}
	hc.header[key] = append(hc.header[key], value)
	return hc
}

func (hc *HttpClient) SetParams(param map[string][]string) *HttpClient {
	hc.param = param
	return hc
}

func (hc *HttpClient) SetParam(key, value string) *HttpClient {
	if hc.param == nil {
		hc.param = map[string][]string{}
	}
	if "" != value {
		hc.param[key] = []string{value}
	} else {
		delete(hc.param, key)
	}
	return hc
}

func (hc *HttpClient) AddParam(key, value string) *HttpClient {
	if hc.param == nil {
		hc.param = map[string][]string{}
	}
	hc.param[key] = append(hc.param[key], value)
	return hc
}

// SetBody 实际POST请求数据，可以上传文件。此参数设置后hc.param失效。
func (hc *HttpClient) SetBody(body any) *HttpClient {
	switch body.(type) {
	case io.Reader:
		hc.body = body.(io.Reader)
	case []byte:
		hc.body = bytes.NewReader(body.([]byte))
	case string:
		hc.body = strings.NewReader(body.(string))
	default:
		jsonBuf, err := json.Marshal(body)
		if err != nil {
			hc.error = errors.Join(hc.error, err)
		} else {
			hc.body = bytes.NewReader(jsonBuf)
		}
	}
	return hc
}

// EncodeUrl url编码
func (hc *HttpClient) EncodeUrl() *HttpClient {
	if strings.ContainsRune(hc.url, '?') {
		u, err := url.Parse(hc.url)
		if err != nil {
			hc.error = errors.Join(hc.error, err)
		} else {
			hc.url = hc.url[:strings.Index(hc.url, "?")+1] + u.Query().Encode()
		}
	}
	return hc
}

// EncodeParam 表单参数编码
func (hc *HttpClient) EncodeParam() *HttpClient {
	if hc.param != nil {
		for _, vs := range hc.param {
			for i := 0; i < len(vs); i++ {
				vs[i] = url.QueryEscape(vs[i])
			}
		}
	}
	return hc
}

func (hc *HttpClient) SetContentType(contentType httpContentType) *HttpClient {
	hc.SetHeader("Content-Type", string(contentType))
	return hc
}

func (hc *HttpClient) SetTimeout(second int64) *HttpClient {
	hc.timeout = second
	return hc
}

// Do 发送http请求
func (hc *HttpClient) Do(method httpMethod) (response *HttpResponse, err error) {
	//返回构建过程中的错误信息
	if hc.error != nil {
		return nil, hc.error
	}
	//组装表单参数
	var paramPair string
	if hc.param != nil {
		var buf bytes.Buffer
		for k, vs := range hc.param {
			for _, v := range vs {
				buf.WriteString("&")
				buf.WriteString(k)
				buf.WriteString("=")
				buf.WriteString(v)
			}
		}
		paramPair = buf.String()
	}
	//创建client
	c := http.Client{Timeout: time.Duration(hc.timeout) * time.Second}
	if hc.transport != nil {
		c.Transport = hc.transport
	}
	//创建request
	var request *http.Request
	switch method {
	case HttpMethodGet: //GET请求
		//在url上附加表单参数
		if paramPair != "" {
			if strings.ContainsRune(hc.url, '?') {
				hc.url += paramPair
			} else {
				hc.url += "?" + paramPair
			}
		}
		request, err = http.NewRequest(http.MethodGet, hc.url, nil)
	case HttpMethodPost: //POST请求
		//未设置body时，默认将表单数据转换为body。
		if paramPair != "" && hc.body == nil {
			hc.body = strings.NewReader(paramPair)
		}
		request, err = http.NewRequest(http.MethodPost, hc.url, hc.body)
	}
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, errors.New("not support http method")
	}
	if hc.header != nil {
		for k, vs := range hc.header {
			for _, v := range vs {
				request.Header.Add(k, v)
			}
		}
	}
	//发送http请求
	resp, err := c.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			fmt.Println("http response close err:", err)
		}
	}(resp.Body)
	//处理response
	response = &HttpResponse{
		StatusCode:    resp.StatusCode,
		Header:        resp.Header,
		ContentLength: resp.ContentLength,
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response.Body = body
	//返回结果
	return response, err
}

// Get GET请求
func (hc *HttpClient) Get() (response *HttpResponse, err error) {
	return hc.Do(HttpMethodGet)
}

// PostForm 普通POST请求：application/x-www-form-urlencoded
func (hc *HttpClient) PostForm() (response *HttpResponse, err error) {
	return hc.SetContentType(HttpContentTypeEncoded).Do(HttpMethodPost)
}

// PostMulti 可以上传文件POST请求：multipart/form-data
func (hc *HttpClient) PostMulti() (response *HttpResponse, err error) {
	return hc.SetContentType(HttpContentTypeForm).Do(HttpMethodPost)
}

// PostJson json格式的POST请求：application/json
func (hc *HttpClient) PostJson() (response *HttpResponse, err error) {
	return hc.SetContentType(HttpContentTypeJson).Do(HttpMethodPost)
}

// PostXml xml格式的POST请求：application/xml
func (hc *HttpClient) PostXml() (response *HttpResponse, err error) {
	return hc.SetContentType(HttpContentTypeXml).Do(HttpMethodPost)
}

func (hr *HttpResponse) ToString() string {
	return string(hr.Body)
}

func (hr *HttpResponse) ToMap() (result map[string]interface{}, err error) {
	result = make(map[string]interface{})
	err = json.Unmarshal(hr.Body, &result)
	return
}

func (hr *HttpResponse) ToAny(t any) error {
	return json.Unmarshal(hr.Body, &t)
}
