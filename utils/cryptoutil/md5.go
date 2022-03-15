//MD5的全称是Message-DigestAlgorithm 5,它可以把一个任意长度的字节数组转换成一个定长的整数,并且这种转换是不可逆的.
//对于任意长度的数据,转换后的MD5值长度是固定的,而且MD5的转换操作很容易,只要原数据有一点点改动,转换后结果就会有很大的差异.
//正是由于MD5算法的这些特性,它经常用于对于一段信息产生信息摘要,以防止其被篡改.其还广泛就于操作系统的登录过程中的安全验证,
//比如Unix操作系统的密码就是经过MD5加密后存储到文件系统中,当用户登录时输入密码后, 对用户输入的数据经过MD5加密后与原来存储的密文信息比对,
//如果相同说明密码正确,否则输入的密码就是错误的. MD5以512位为一个计算单位对数据进行分组,每一分组又被划分为16个32位的小组,经过一系列处理后,
//输出4个32位的小组,最后组成一个128位的哈希值.对处理的数据进行512求余得到N和一个余数,如果余数不为448,填充1和若干个0直到448位为止,
//最后再加上一个64位用来保存数据的长度,这样经过预处理后,数据变成（N+1）x 512位.

package cryptoutil

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

// MD5 简单数据md5加密(不适用于并发环境)
func MD5(s string) string {
	b := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", b)
}

// MD5Check 函数传入一个未加密的字符串和与加密后的数据,进行对比,如果正确就返回true
func MD5Check(content, encrypted string) bool {
	return strings.EqualFold(MD5Encode(content), encrypted)
}

// MD5Encode 函数用来加密数据
func MD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
