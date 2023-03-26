package netutil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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
