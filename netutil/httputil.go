package netutil

import (
	"bytes"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

type HttpClient struct {
	url    string
	header map[string][]string
	param  map[string][]string
	body   io.Reader
}

type HttpResponse struct {
	StatusCode    int // e.g. 200
	Header        map[string][]string
	Body          []byte
	ContentLength int64
}

func NewHttpClient(url string) *HttpClient {
	return &HttpClient{
		url: url,
	}
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
			fmt.Println("SetBody error：" + err.Error())
		}
		hc.body = bytes.NewReader(jsonBuf)
	}
	return hc
}

func (hc *HttpClient) EncodeUrl() *HttpClient {
	if strings.ContainsRune(hc.url, '?') {
		u, _ := url.Parse(hc.url)
		hc.url = u.Path + "?" + u.Query().Encode()
	}
	return hc
}

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

func (hc *HttpClient) Do(method httpMethod) (response *HttpResponse, err error) {
	c := http.Client{Timeout: 30 * time.Second}
	var request *http.Request
	switch method {
	case HttpMethodGet:
		apiUrl := hc.url
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
			if strings.ContainsRune(apiUrl, '?') {
				apiUrl += buf.String()
			} else {
				apiUrl += "?" + buf.String()
			}
		}
		request, err = http.NewRequest(http.MethodGet, apiUrl, nil)
	case HttpMethodPost:
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
	return response, err
}

func (hc *HttpClient) Get() (response *HttpResponse, err error) {
	return hc.Do(HttpMethodGet)
}

func (hc *HttpClient) Post() (response *HttpResponse, err error) {
	return hc.SetContentType(HttpContentTypeEncoded).Do(HttpMethodPost)
}

func (hc *HttpClient) PostForm() (response *HttpResponse, err error) {
	return hc.SetContentType(HttpContentTypeForm).Do(HttpMethodPost)
}

func (hc *HttpClient) PostJson() (response *HttpResponse, err error) {
	return hc.SetContentType(HttpContentTypeJson).Do(HttpMethodPost)
}

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

// Get Get请求
func Get(url string, param map[string]string) ([]byte, error) {
	if param != nil {
		var buf bytes.Buffer
		for k, v := range param {
			buf.WriteString("&")
			buf.WriteString(k)
			buf.WriteString("=")
			buf.WriteString(v)
		}
		if strings.ContainsRune(url, '?') {
			url += buf.String()
		} else {
			url += "?" + buf.String()
		}
	}
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return ioutil.ReadAll(rsp.Body)
}

// GetWithEncode Get请求，参数编码
func GetWithEncode(url string, param map[string]string) ([]byte, error) {
	if param != nil {
		if strings.ContainsRune(url, '?') {
			url += urlEncode(param)
		} else {
			url = url + "?" + urlEncode(param)
		}
	}
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return ioutil.ReadAll(rsp.Body)
}

func urlEncode(param map[string]string) string {
	var buf bytes.Buffer
	for k, v := range param {
		buf.WriteString("&")
		buf.WriteString(url.QueryEscape(k))
		buf.WriteString("=")
		buf.WriteString(url.QueryEscape(v))
	}
	return buf.String()
}

// PostForm 简单的POST方法
func PostForm(apiUrl string, param url.Values) (map[string]interface{}, error) {
	resp, err := http.PostForm(apiUrl, param)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	json.Unmarshal(b, &result)
	return result, nil
}

// PostFormWidthHeader 表单格式的POST方法
func PostFormWidthHeader(apiUrl string, header map[string]string, param url.Values) (map[string]interface{}, error) {
	c := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(param.Encode()))
	if err != nil {
		return nil, err
	}
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

// PostJson JSON格式的POST方法
func PostJson(apiUrl string, header map[string]string, param map[string]interface{}) (map[string]interface{}, error) {
	p, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	c := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(p))
	if err != nil {
		return nil, err
	}
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

// PostFile 带文件上传的POST方法
func PostFile() {

}
