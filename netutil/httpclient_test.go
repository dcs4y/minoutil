package netutil

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

// 证书示例
func Test_Transport(t *testing.T) {
	caCert, err := os.ReadFile("path/to/ca.crt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	cert, err := tls.LoadX509KeyPair("path/to/client.crt", "path/to/client.key")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{cert},
		},
	}
	resp, err := NewHttpClient("https://www.baidu.com").SetTransport(transport).Get()
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(resp)
	}
}

// 上传文件并提交表单参数示例
func Test_FileUpload(t *testing.T) {
	url := "http://127.0.0.1:8888/upload"
	filePath := "example.zip"

	// 创建 multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件字段
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		fmt.Println("创建表单字段失败:", err)
		return
	}
	if _, err = io.Copy(part, file); err != nil {
		fmt.Println("写入文件失败:", err)
		return
	}
	// 添加其他字段（可选）
	err = writer.WriteField("key", "value")
	if err != nil {
		fmt.Println("添加其他字段失败:", err)
		return
	}
	err = writer.Close()
	if err != nil {
		return
	}
	resp, err := NewHttpClient(url).SetBody(body).SetHeader("Content-Type", writer.FormDataContentType()).Do(HttpMethodPost)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(resp)
	}
}
