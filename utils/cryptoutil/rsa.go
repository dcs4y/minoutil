//首先使用openssl生成公私钥,使用RSA的时候需要提供公钥和私钥 , 可以通过openssl来生成对应的pem格式的公钥和私钥匙

package cryptoutil

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"os"
)

// KeyPairs 使用golang标准库ecdsa生成非对称(ES256,ES384,ES521)加密密钥对
func KeyPairs(size int, keyName string) {
	var c elliptic.Curve
	switch size {
	case 256:
		c = elliptic.P256()
	case 384:
		c = elliptic.P384()
	case 521:
		c = elliptic.P521()
	}
	privateKey, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	privateBs := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	privateFile, err := os.Create(keyName + ".private.pem")
	if err != nil {
		log.Fatal(err)
	}
	_, err = privateFile.Write(privateBs)
	if err != nil {
		log.Fatal(err)
	}
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(privateKey.Public())
	publicBs := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	publicKeyFile, err := os.Create(keyName + ".public.pem")
	if err != nil {
		log.Fatal(err)
	}
	_, err = publicKeyFile.Write(publicBs)
	if err != nil {
		log.Fatal(err)
	}
}

// RsaEncrypt 加密
func RsaEncrypt(origData []byte, publicKey []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// RsaDecrypt 解密
func RsaDecrypt(ciphertext []byte, privateKey []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
