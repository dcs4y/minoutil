//当前我们项目中常用的加解密的方式无非三种.
//对称加密, 加解密都使用的是同一个密钥, 其中的代表就是AES,DES
//非对加解密, 加解密使用不同的密钥, 其中的代表就是RSA
//签名算法, 如MD5,SHA1,HMAC等, 主要用于验证,防止信息被修改, 如：文件校验,数字签名,鉴权协议

package cryptoutil

import (
	"encoding/base64"
	"fmt"
	"log"
	"testing"
)

func Test_aes(t *testing.T) {
	orig := "hello world"
	key := "123456781234567812345678"
	fmt.Println("原文：", orig)

	encryptCode := AesEncrypt(orig, key)
	fmt.Println("密文：", encryptCode)

	decryptCode := AesDecrypt(encryptCode, key)
	fmt.Println("解密结果：", decryptCode)
}

func Test_base64(t *testing.T) {
	b := NewBase64("IJjkKLMNO567PQX12RVW3YZaDEFGbcdefghiABCHlSTUmnopqrxyz04stuvw89+/")
	d := b.Encode([]byte("它是一种数据编码方式,虽然是可逆的,但是它的编码方式是公开的,无所谓加密."))
	fmt.Println(string(d))
	d, err := b.Decode(d)
	fmt.Println(string(d), err)
}

func Test_des(t *testing.T) {
	key := []byte("2fa6c1e9")
	str := "I love this beautiful world!"
	strEncrypted, err := DesEncrypt(str, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Encrypted:", strEncrypted)
	strDecrypted, err := DesDecrypt(strEncrypted, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Decrypted:", strDecrypted)
}

func Test_md5(t *testing.T) {
	s := "aa6d11692a0e11eb92f800163e0ad1ee" + "66b39c5b3abb11eca4f600163e0ad1ee" + "11"
	fmt.Println(MD5(s))
	fmt.Println(MD5Encode(s))
}

func Test_md5_check(t *testing.T) {
	strTest := "I love this beautiful world!"
	strEncrypted := "98b4fc4538115c4980a8b859ff3d27e1"
	fmt.Println(MD5Check(strTest, strEncrypted))
}

func Test_keyPairs(t *testing.T) {
	KeyPairs(521, "keys")
}

func Test_rsa(t *testing.T) {
	// 私钥生成
	//openssl genrsa -out rsa_private_key.pem 1024
	var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQDcGsUIIAINHfRTdMmgGwLrjzfMNSrtgIf4EGsNaYwmC1GjF/bM
h0Mcm10oLhNrKNYCTTQVGGIxuc5heKd1gOzb7bdTnCDPPZ7oV7p1B9Pud+6zPaco
qDz2M24vHFWYY2FbIIJh8fHhKcfXNXOLovdVBE7Zy682X1+R1lRK8D+vmQIDAQAB
AoGAeWAZvz1HZExca5k/hpbeqV+0+VtobMgwMs96+U53BpO/VRzl8Cu3CpNyb7HY
64L9YQ+J5QgpPhqkgIO0dMu/0RIXsmhvr2gcxmKObcqT3JQ6S4rjHTln49I2sYTz
7JEH4TcplKjSjHyq5MhHfA+CV2/AB2BO6G8limu7SheXuvECQQDwOpZrZDeTOOBk
z1vercawd+J9ll/FZYttnrWYTI1sSF1sNfZ7dUXPyYPQFZ0LQ1bhZGmWBZ6a6wd9
R+PKlmJvAkEA6o32c/WEXxW2zeh18sOO4wqUiBYq3L3hFObhcsUAY8jfykQefW8q
yPuuL02jLIajFWd0itjvIrzWnVmoUuXydwJAXGLrvllIVkIlah+lATprkypH3Gyc
YFnxCTNkOzIVoXMjGp6WMFylgIfLPZdSUiaPnxby1FNM7987fh7Lp/m12QJAK9iL
2JNtwkSR3p305oOuAz0oFORn8MnB+KFMRaMT9pNHWk0vke0lB1sc7ZTKyvkEJW0o
eQgic9DvIYzwDUcU8wJAIkKROzuzLi9AvLnLUrSdI6998lmeYO9x7pwZPukz3era
zncjRK3pbVkv0KrKfczuJiRlZ7dUzVO0b6QJr8TRAA==
-----END RSA PRIVATE KEY-----
`)

	// 公钥: 根据私钥生成
	//openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
	var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDcGsUIIAINHfRTdMmgGwLrjzfM
NSrtgIf4EGsNaYwmC1GjF/bMh0Mcm10oLhNrKNYCTTQVGGIxuc5heKd1gOzb7bdT
nCDPPZ7oV7p1B9Pud+6zPacoqDz2M24vHFWYY2FbIIJh8fHhKcfXNXOLovdVBE7Z
y682X1+R1lRK8D+vmQIDAQAB
-----END PUBLIC KEY-----
`)

	data, _ := RsaEncrypt([]byte("hello world"), publicKey)
	fmt.Println(base64.StdEncoding.EncodeToString(data))
	origData, _ := RsaDecrypt(data, privateKey)
	fmt.Println(string(origData))
}

func Test_sha1(t *testing.T) {
	fmt.Println(Sha1("hello world!"))
}

func Test_sha256(t *testing.T) {
	fmt.Println(Sha256("hello world!", "password"))
}
