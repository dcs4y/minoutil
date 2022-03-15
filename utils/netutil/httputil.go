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

// PostSample 简单的POST方法
func PostSample(apiUrl string, param url.Values) (map[string]interface{}, error) {
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

// PostForm 表单格式的POST方法
func PostForm(apiUrl string, header map[string]string, param url.Values) (map[string]interface{}, error) {
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
