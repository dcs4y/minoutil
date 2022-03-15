// hashid 将有序数字转化为无序且唯一的字符串

package hashidutil

import (
	"github.com/speps/go-hashids/v2"
)

var hd = NewHashIDData()

func NewHashIDData() *hashids.HashIDData {
	hd := hashids.NewData()
	hd.Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	hd.Salt = "DoYouKnowThisIsHashIdSalt"
	hd.MinLength = 4
	return hd
}

// EncodeHashId 数字转字符串
func EncodeHashId(n int64) string {
	h, _ := hashids.NewWithData(hd)
	e, _ := h.EncodeInt64([]int64{n})
	return e
}

func EncodeHashHex(hex string) string {
	h, _ := hashids.NewWithData(hd)
	e, _ := h.EncodeHex(hex)
	return e
}

// DecodeHashId 字符串转数字
func DecodeHashId(s string) (int64, error) {
	h, _ := hashids.NewWithData(hd)
	d, err := h.DecodeInt64WithError(s)
	if err != nil {
		return 0, err
	} else {
		return d[0], nil
	}
}

func DecodeHashHex(s string) (string, error) {
	h, _ := hashids.NewWithData(hd)
	return h.DecodeHex(s)
}
