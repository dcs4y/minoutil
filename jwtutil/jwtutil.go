package jwtutil

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// https://github.com/dgrijalva/jwt-go
type JWTClaims struct {
	jwt.StandardClaims
	UserId   int64  `json:"userId"`
	UserName string `json:"userName"`
	UserType int    `json:"userType"`
	RoleId   string `json:"roleId"`
}

// GenToken generate jwt token
func GenToken(claims *JWTClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("token生成失败！")
	}
	return signedToken, nil
}

// VerifyAction 验证jwt token
func VerifyAction(strToken string, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.New("token解析失败！")
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("token转换失败！")
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, errors.New("token验证失败！")
	}
	return claims, nil
}

// GetClaims 从上下文中获取登录信息
func GetClaims(c *gin.Context) *JWTClaims {
	v, _ := c.Get("claims")
	return v.(*JWTClaims)
}

func GetClaimsWithCheck(c *gin.Context) (*JWTClaims, bool) {
	v, b := c.Get("claims")
	if b {
		return v.(*JWTClaims), b
	} else {
		return nil, b
	}
}
