package qrcodeutil

import (
	qrEncode "github.com/skip2/go-qrcode"
)

// https://github.com/skip2/go-qrcode

// Encode 生成二维码内容
func Encode(content string, size int) ([]byte, error) {
	return qrEncode.Encode(content, qrEncode.Medium, size)
}

// EncodeToFile 生成二维码到文件
func EncodeToFile(content string, size int, filename string) error {
	return qrEncode.WriteFile(content, qrEncode.Medium, size, filename)
}
