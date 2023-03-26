package qrcodeutil

import (
	qrDecode "github.com/tuotoo/qrcode"
	"os"
)

// https://github.com/tuotoo/qrcode

// Decode 识别二维码
func Decode(filename string) (string, error) {
	fi, err := os.Open(filename)
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
