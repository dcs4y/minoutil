package jwtutil

import (
	"crypto/rsa"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RequestContextTokenName 接口请求头中登录标识的名称
var RequestContextTokenName = "token"

// ContextJWTClaims 解析后的登录信息名称
var ContextJWTClaims = "claims"

// https://github.com/golang-jwt/jwt
type JWTClaims struct {
	jwt.MapClaims
	ExpiresAt  int64  `json:"expiresAt"` //有效期。Unix时间(秒)。
	UserId     string `json:"userId"`
	UserName   string `json:"userName"`
	UserType   int    `json:"userType"`
	ClientType int    `json:"clientType"`
	RoleId     string `json:"roleId"`
}

// GenerateToken generate jwt token
// 对称加密（HS256）
func GenerateToken(claims JWTClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// VerifyToken 验证jwt token
func VerifyToken(strToken string, secret string) (JWTClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	var claims JWTClaims
	if err != nil || !token.Valid {
		return claims, errors.New("token验证失败！")
	}
	claims, ok := token.Claims.(JWTClaims)
	if !ok {
		return claims, errors.New("token转换失败！")
	}
	return claims, nil
}

// GenerateRS256Token 使用非对称加密（RS256）
func GenerateRS256Token(claims JWTClaims, privateKey *rsa.PrivateKey) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func VerifyRS256Token(tokenString string, publicKey *rsa.PublicKey) (JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	var claims JWTClaims
	if err != nil || !token.Valid {
		return claims, errors.New("token验证失败！")
	}
	claims, ok := token.Claims.(JWTClaims)
	if !ok {
		return claims, errors.New("token转换失败！")
	}
	return claims, nil
}

// GetClaims 从上下文中获取登录信息
func GetClaims(c *gin.Context) JWTClaims {
	v, _ := c.Get(ContextJWTClaims)
	return v.(JWTClaims)
}

func GetClaimsWithCheck(c *gin.Context) (claims JWTClaims, b bool) {
	v, b := c.Get(ContextJWTClaims)
	if b {
		return v.(JWTClaims), b
	} else {
		return claims, b
	}
}
