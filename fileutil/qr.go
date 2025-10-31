package fileutil

import (
	"io"
	"os"

	qrEncode "github.com/skip2/go-qrcode"
	qrDecode "github.com/tuotoo/qrcode"
)

// https://github.com/skip2/go-qrcode

// QREncode 生成二维码内容
func QREncode(content string, size int) ([]byte, error) {
	return qrEncode.Encode(content, qrEncode.Medium, size)
}

// QREncodeToFile 生成二维码到文件
func QREncodeToFile(content string, size int, filename string) error {
	return qrEncode.WriteFile(content, qrEncode.Medium, size, filename)
}

// https://github.com/tuotoo/qrcode

// QRDecode 识别二维码
func QRDecode(r io.Reader) (string, error) {
	qr, err := qrDecode.Decode(r)
	if err != nil {
		return "", err
	}
	return qr.Content, nil
}

// QRDecodeFromFile 从文件识别二维码
func QRDecodeFromFile(filePath string) (string, error) {
	fi, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer fi.Close()
	qr, err := qrDecode.Decode(fi)
	if err != nil {
		return "", err
	}
	return qr.Content, nil
}
