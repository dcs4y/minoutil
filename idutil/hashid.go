package idutil

import (
	"github.com/speps/go-hashids/v2"
)

// HashId 将有序数字转化为无序且唯一的字符串
var HashId = NewHashIDData(4)

type hashId struct {
	*hashids.HashIDData
}

func NewHashIDData(minLength int) *hashId {
	hd := hashId{hashids.NewData()}
	hd.Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	hd.Salt = "DoYouKnowThisIsHashIdSalt"
	hd.MinLength = minLength
	return &hd
}

// EncodeHashId 数字转字符串
func (hd *hashId) EncodeHashId(n int64) string {
	h, _ := hashids.NewWithData(hd.HashIDData)
	e, _ := h.EncodeInt64([]int64{n})
	return e
}

// EncodeHashHex 16进制转字符串
func (hd *hashId) EncodeHashHex(hex string) string {
	h, _ := hashids.NewWithData(hd.HashIDData)
	e, _ := h.EncodeHex(hex)
	return e
}

// DecodeHashId 字符串转数字
func (hd *hashId) DecodeHashId(s string) (int64, error) {
	h, _ := hashids.NewWithData(hd.HashIDData)
	d, err := h.DecodeInt64WithError(s)
	if err != nil {
		return 0, err
	} else {
		return d[0], nil
	}
}

// DecodeHashHex 字符串转16进制
func (hd *hashId) DecodeHashHex(s string) (string, error) {
	h, _ := hashids.NewWithData(hd.HashIDData)
	return h.DecodeHex(s)
}

// EncodeHashId 数字转字符串
func EncodeHashId(n int64) string {
	return HashId.EncodeHashId(n)
}

// EncodeHashHex 16进制转字符串
func EncodeHashHex(hex string) string {
	return HashId.EncodeHashHex(hex)
}

// DecodeHashId 字符串转数字
func DecodeHashId(s string) (int64, error) {
	return HashId.DecodeHashId(s)
}

// DecodeHashHex 字符串转16进制
func DecodeHashHex(s string) (string, error) {
	return HashId.DecodeHashHex(s)
}
