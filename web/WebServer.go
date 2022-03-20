package web

import (
	"context"
	"game/common"
	"game/utils/logutil"
	"game/web/internal/handler"
	"game/web/internal/middleware"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"
)

var log = logutil.Log

// Start WEB服务启动入口
// 示例：模板，登录接口，跨域，jwt，鉴权，用户，角色，功能权限，文件上传下载，websocket，swagger，导入导出，生成pdf。
// https://github.com/gin-gonic/gin
func Start() {
	if !common.WS.WebConfig.Enable {
		log.Println("WEB服务未运行！")
		return
	}
	g := gin.Default()

	// 限流
	g.Use(middleware.Limiter())
	// 日志
	g.Use(middleware.LoggerToFile())
	// 跨域
	g.Use(middleware.Cors())

	// 8 MiB 设置最大的上传文件的大小
	g.MaxMultipartMemory = 8 << 20

	// 静态资源
	g.StaticFS("image", http.Dir(filepath.Join(common.ResourcesPath, "image")))
	g.StaticFS("game", http.Dir(filepath.Join(common.ResourcesPath, "page")))

	// 模板设置
	g.SetFuncMap(template.FuncMap{
		"unescaped": func(x string) template.HTML {
			return template.HTML(x)
		},
	})
	g.LoadHTMLGlob(filepath.Join(common.ResourcesPath, "templates/*"))

	// 公共访问业务
	handler.HandleWithPublic(g)
	// jwt登录认证
	g.Use(middleware.JWTAuth())
	// 仅需要登录的业务
	handler.HandleWithLogin(g)
	// 鉴权
	g.Use(middleware.RBACAuth())
	// 处理需要鉴权业务
	handler.HandleWithAuth(g)

	log.Println("WEB服务监听：", common.WS.WebConfig.Port)

	// 启动服务
	//g.Run(":" + strconv.Itoa(common.WS.WebConfig.Port))

	// 优雅重启或停止
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(common.WS.WebConfig.Port),
		Handler: g,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen :%s \n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server...")
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
