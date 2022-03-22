//当前我们项目中常用的加解密的方式无非三种.
//对称加密, 加解密都使用的是同一个密钥, 其中的代表就是AES,DES
//非对加解密, 加解密使用不同的密钥, 其中的代表就是RSA
//签名算法, 如MD5,SHA1,HMAC等, 主要用于验证,防止信息被修改, 如：文件校验,数字签名,鉴权协议

package cryptoutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
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
	prvkey, pubkey := RsaGenerateKey(2048)
	fmt.Println("私钥：", string(prvkey))
	fmt.Println("公钥：", string(pubkey))
}

func Test_rsa(t *testing.T) {
	// 私钥生成
	//openssl genrsa -out rsa_private_key.pem 1024
	var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAu86O3rCE2MdzQLyzQUmwp7aWCHOmCKu/JgSufk7LRqsV3ein
6Kqxgk5350IhTTvhIW/MvxTRRR20xrfW965hPwtpIS6dm8ChDcrjd97p24lXYAaz
rQ2aIxuZiHG4+scPqRvmfb6jyuQ22PwdtGtJjhXj/Z0INCuNmBhmPDKgF3N96xca
muJYqVmIiD3BDslqHxFKq4N41ECWObgN3tuLm4nHLCkf2aBxDmPPWZQeLEcCUQZg
WTgx0ix7/XVaH86X4QwvgwzbcW57bKDxLCvflskvoqHnmJCH50ZLOBPiEJWRQ5R8
LIneF1Vgl2ytTMvxzscBuUgisBaAP6WsjgraswIDAQABAoIBAD39cCsRGMh1DRXR
M1nZePXizqL7iVJTXkSuRupqF667yfv1T3b84JqiS/GJYnSbzzO6M1rfBDRMGd99
zvbyGCc3HPxW5q8CZianUW2/pnFQZAbOL4BvfPEZqxPedbBRBFpNW0cmJepSacg1
b5id0SmVECwmKQ8PUS1i4Fv+WdljNARYSmVbTSVb9M3O8yYkk6FM9jdyJOcw2uu9
vpk+na1dLTxoogtnHKBk4sbPVnzrAV/P7HYfpgVjvZi0n1h7+ROjqDAFXRdHpE7H
7qKuDZ0S/TTFAFCjJ6X2h5GXbc3MzaMmLW8Q9BjkRWJ3jttszM+nDeUUtov1V99b
Hv8LEAECgYEAxA5sAbLjZl9VCFDBVVfUwWe0aBCNpKlSf5z2exr0HtWwzZoppaNY
WZd+nLnTdsYYQY5r5+OMRsGPIHcVQmXo747tBbAElMLW37jJZgtH7/R0UtVjDQ2F
tVTx6q+3Wt3Ok9JFs1AZmsk2ysghn7IT1SUA9h2sgylFUTPab1EFFZECgYEA9Tpw
mnLpJET8ngSyL0CFZT0JSxLonhe81SNH44jcK9+o8bGMwAi6vgC5OzIYKkigYfUk
7cjXg57jwpKoI7B8K+xKYmVtnK+sHDHNulcVWkuaO2JtLx2Os2kSs8AKvIhtSuMP
uXY1hvnClPj8dSSd4mK40ebwH+W6y5UJXX8s+gMCgYBW/zNMGa4wahMYaoUvspa0
76itGNNRgtUZzXPOMqqq2BXpVgQu/Omib8f+EbNVHBf9Vw5oyp8fcpppRI5JdWFE
k/53LKELxd2FTsEHp+/W+Xl4nDmkvCBd04C5rBlHl+8nxwGozN1fwvOuTdolu+Zz
CUWA8K/xT6nzm3fNN31zEQKBgQCZ6K+7f9trdBDIbFIY7fnK7F/kyl1Mu8E8VARO
fhsRPjKeXrzj2Za6oWGgEXSxNRvT3zPqOV9psFNqTvlQTPfh6U9WKip9aJQwreFc
zkMMhN6r3r6AD6D8YTnnruQOJ+HQWmoIEHTP3fmN3ic0rmKZLzSoKLUUj+6Iz2Ut
VSA5WQKBgHb4yZ+AoryTKX+DhELUUUAmfPy974LQfjFgDHKrxW9C18aCPpntxR3U
cm2wtauoU5arHVHhnUZ5FeuhHnYHnMSZQS+7MyjRIHWldfaYySo6fw3wWOnyf69C
KGnbPiX/uzzxnfpffmiBCIHvPVudmskNhilg52ij4nAxJOh78b/2
-----END RSA PRIVATE KEY-----
`)

	// 公钥: 根据私钥生成
	//openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
	var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu86O3rCE2MdzQLyzQUmw
p7aWCHOmCKu/JgSufk7LRqsV3ein6Kqxgk5350IhTTvhIW/MvxTRRR20xrfW965h
PwtpIS6dm8ChDcrjd97p24lXYAazrQ2aIxuZiHG4+scPqRvmfb6jyuQ22PwdtGtJ
jhXj/Z0INCuNmBhmPDKgF3N96xcamuJYqVmIiD3BDslqHxFKq4N41ECWObgN3tuL
m4nHLCkf2aBxDmPPWZQeLEcCUQZgWTgx0ix7/XVaH86X4QwvgwzbcW57bKDxLCvf
lskvoqHnmJCH50ZLOBPiEJWRQ5R8LIneF1Vgl2ytTMvxzscBuUgisBaAP6Wsjgra
swIDAQAB
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

func Test_ke(t *testing.T) {
	prvKey, pubKey := GenRsaKey()
	fmt.Println(string(prvKey))
	fmt.Println(string(pubKey))
}

//RSA公钥私钥产生
func GenRsaKey() (prvkey, pubkey []byte) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvkey = pem.EncodeToMemory(block)
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey = pem.EncodeToMemory(block)
	return
}
