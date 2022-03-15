package middleware

import (
	"game/common"
	"game/utils/limiterutil"
	"github.com/gin-gonic/gin"
	"net/http"
)

var pathMap = make(map[string]string)

// Limiter uri级别限流
func Limiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if p, b := pathMap[path]; b {
			l := limiterutil.NewLimiter(p, 10, 10)
			if !l.Allow() {
				// 被限制访问
				c.JSON(http.StatusOK, common.ResultErrorMessage("服务器繁忙，请等待。。。"))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
