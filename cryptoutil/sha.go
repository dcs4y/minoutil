package cryptoutil

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Sha1 算法
func Sha1(s string) string {
	//产生一个散列值得方式是 sha1.New(),sha1.Write(bytes),然后 sha1.Sum([]byte{}).这里我们从一个新的散列开始.
	h := sha1.New()
	//写入要处理的字节.如果是一个字符串,需要使用[]byte(s) 来强制转换成字节数组.
	h.Write([]byte(s))
	//这个用来得到最终的散列值的字符切片.Sum 的参数可以用来都现有的字符切片追加额外的字节切片：一般不需要要.
	bs := h.Sum(nil)
	//SHA1 值经常以 16 进制输出,例如在 git commit 中.使用%x 来将散列结果格式化为 16 进制字符串.
	return fmt.Sprintf("%x", bs)
}

// Sha256 算法
func Sha256(data, secret string) string {
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))
	// Write Data to it
	h.Write([]byte(data))
	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
