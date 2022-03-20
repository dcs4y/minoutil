package handler

import (
	"game/common"
	"game/entity/param"
	"game/entity/result"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

// 获取图形验证码
func captcha(c *gin.Context) {
	// 字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	if id, b64s, err := cp.Generate(); err != nil {
		c.JSON(http.StatusOK, common.ResultErrorMessage("验证码获取失败"))
	} else {
		c.JSON(http.StatusOK, common.ResultOkData(result.CaptchaResponse{
			CaptchaId:     id,
			PicPath:       b64s,
			CaptchaLength: 6,
		}))
	}
}

// 登录
func login(c *gin.Context) {
	var p param.LoginRequest
	c.ShouldBindJSON(&p)
	if err := p.Verify(); err != nil {
		c.JSON(http.StatusOK, common.ResultErrorMessage(err.Error()))
		return
	}
	if store.Verify(p.CaptchaId, p.Captcha, true) {
		// 账号验证
		// 授权token
		c.JSON(http.StatusOK, common.ResultOkData(result.LoginResponse{
			Username: "dcs",
			Token:    "long",
		}))
	} else {
		c.JSON(http.StatusOK, common.ResultErrorMessage("验证码错误"))
	}
}
