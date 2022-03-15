package middleware

import (
	"errors"
	"game/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

// https://github.com/dgrijalva/jwt-go
type JWTClaims struct {
	jwt.StandardClaims
	UserID   int    `json:"userId"`
	UserName string `json:"userName"`
}

var (
	Secret     = "cc.dongs"   //salt
	ExpireTime = 24 * 60 * 60 //token expire time
)

//generate jwt token
func GenToken(claims *JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", errors.New("token生成失败！")
	}
	return signedToken, nil
}

//验证jwt token
func verifyAction(strToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
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

// JWTAuth 定义一个JWTAuth的中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 通过http header中的token解析来认证
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusOK, common.ResultErrorMessage("非法访问！"))
			c.Abort()
			return
		}
		claims, err := verifyAction(token)
		if err != nil {
			c.JSON(http.StatusOK, common.ResultErrorMessage(err.Error()))
			c.Abort()
			return
		}
		// 将解析后的有效载荷claims重新写入gin.Context引用对象中
		c.Set("claims", claims)
		// 后续请求处理
		c.Next()
	}
}

// 从上下文中获取登录信息
func GetClaims(c *gin.Context) *JWTClaims {
	v, _ := c.Get("claims")
	return v.(*JWTClaims)
}
