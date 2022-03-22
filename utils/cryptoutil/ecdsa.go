package cryptoutil

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

// EcdsaKeyPairs 使用golang标准库ecdsa生成非对称(ES256,ES384,ES521)加密密钥对
func EcdsaKeyPairs(size int, keyName string) {
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
