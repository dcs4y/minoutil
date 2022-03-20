package handler

import (
	"game/utils/logutil"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	log = logutil.Log.WithField("kind", "web")
	//bootUrl = "https://oapi.dingtalk.com/robot/send?access_token=ca6441d14175831d2f1e4e5409421b7a1a7859824c2cbb59cd6a705f1f318928"
	//secret  = "SEC8864fda8960e963bbac681d27aab862957666aa97b5f0e449e9082a6fd6b26ea"
	//robot   = dingdingclient.NewWebhook(bootUrl, secret)
)

func init() {
	//err := robot.Send(nil, &dingdingclient.TextBody{Text: "WEB服务启动了！"})
	//if err != nil {
	//	log.Println(err)
	//}
}

// HandleWithPublic 公共访问业务
func HandleWithPublic(g *gin.Engine) {
	// websocket。内部处理权限。
	g.GET("ws", GinWebSocketHandler(WsConnHandler))

	// 重定向到首页
	g.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/game/index.html")
	})

	// 图片验证码
	g.POST("/base/captcha", captcha)
	// 登录
	g.POST("/base/login", login)
}

// HandleWithLogin 仅需要登录的业务
func HandleWithLogin(g *gin.Engine) {
	// 获取登录菜单
	g.POST("/menu/getMenu", getMenu)
}

// HandleWithAuth 需要配置相应权限的业务
func HandleWithAuth(g *gin.Engine) {
	// 增加短网址
	//g.POST("/shortUrl/add", AddShortUrl)
}
