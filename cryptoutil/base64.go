//它是一种数据编码方式,虽然是可逆的,但是它的编码方式是公开的,无所谓加密.
//Base64是一种任意二进制到文本字符串的编码方法,常用于在URL,Cookie,网页中传输少量二进制数据. 首先使用Base64编码需要一个含有64个字符的表,
//这个表由大小写字母,数字,+和/组成.采用Base64编码处理数据时,会把每三个字节共24位作为一个处理单元,再分为四组,每组6位,查表后获得相应的字符即编码后的字符串.
//编码后的字符串长32位,这样,经Base64编码后,原字符串增长1/3.如果要编码的数据不是3的倍数,最后会剩下一到两个字节,Base64编码中会采用\x00在处理单元后补全,
//编码后的字符串最后会加上一到两个 = 表示补了几个字节.

package cryptoutil

import (
	"encoding/base64"
)

type base64tool struct {
	coder *base64.Encoding
}

func NewBase64(base64Table string) *base64tool {
	return &base64tool{
		coder: base64.NewEncoding(base64Table),
	}
}

// Encode 编码
func (bt *base64tool) Encode(src []byte) []byte {
	return []byte(bt.coder.EncodeToString(src))
}

// Decode 解码
func (bt *base64tool) Decode(src []byte) ([]byte, error) {
	return bt.coder.DecodeString(string(src))
}
