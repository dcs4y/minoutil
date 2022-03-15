package handler

import (
	"fmt"
	"game/common"
	"game/utils/dingdingclient"
	"game/utils/logutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

var (
	log     = logutil.Log.WithField("kind", "web")
	bootUrl = "https://oapi.dingtalk.com/robot/send?access_token=ca6441d14175831d2f1e4e5409421b7a1a7859824c2cbb59cd6a705f1f318928"
	secret  = "SEC8864fda8960e963bbac681d27aab862957666aa97b5f0e449e9082a6fd6b26ea"
	robot   = dingdingclient.NewWebhook(bootUrl, secret)
)

func init() {
	err := robot.Send(nil, &dingdingclient.TextBody{Text: "WEB服务启动了！"})
	if err != nil {
		log.Println(err)
	}
}

// HandleWithOutAuth 处理无需鉴权业务
func HandleWithOutAuth(g *gin.Engine) {
	// websocket
	g.GET("ws", GinWebSocketHandler(WsConnHandler))

	// 重定向到首页
	g.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/game/home.html")
	})

	// 首页
	//g.GET("tmpl/home", homePage)
	//g.GET("api/home", home)

	// 单文件上传
	g.GET("tmpl/uploadFile", func(context *gin.Context) {
		context.HTML(http.StatusOK, "uploadfile.html", nil)
	})
	g.POST("api/uploadFile", func(context *gin.Context) {
		f, err := context.FormFile("f1")
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		context.SaveUploadedFile(f, "static/"+f.Filename)
		context.JSON(http.StatusOK, gin.H{"meesage": "OK"})
	})

	// 多文件上传
	g.GET("tmpl/uploadFiles", func(context *gin.Context) {
		filePathNameArray, err := filepath.Glob(filepath.Join(common.ResourcesPath, "data/files", "*.jpg"))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
		for i, path := range filePathNameArray {
			_, file := filepath.Split(path)
			filePathNameArray[i] = file
		}
		context.HTML(http.StatusOK, "uploadfiles.html", filePathNameArray)
	})
	g.POST("api/uploadFiles", func(context *gin.Context) {
		form, err := context.MultipartForm()
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		files := form.File["f1s"]
		for _, f := range files {
			fmt.Println(f.Filename)
			context.SaveUploadedFile(f, "static/"+f.Filename)
		}
		//context.JSON(http.StatusOK, gin.H{"message": "OK"})
		context.Redirect(http.StatusMovedPermanently, "uploadFiles")
	})
}

// HandleWithAuth 处理需要鉴权业务
func HandleWithAuth(g *gin.Engine) {
	// 增加短网址
	//g.POST("/shortUrl/add", AddShortUrl)
}
