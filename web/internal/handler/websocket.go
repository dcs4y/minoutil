package handler

import (
	"game/common"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"time"
)

// GinWebSocketHandler websocket.Handler 转 gin HandlerFunc
func GinWebSocketHandler(wsConnHandler websocket.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("new ws request:%v\n", c.Request.RemoteAddr)
		if c.IsWebsocket() {
			wsConnHandler.ServeHTTP(c.Writer, c.Request)
		} else {
			_, _ = c.Writer.WriteString("=== not websocket request ===")
		}
	}
}

// WsConnHandler websocket连接处理
func WsConnHandler(conn *websocket.Conn) {
	for {
		var message string
		if err := websocket.Message.Receive(conn, &message); err != nil {
			log.Println(err)
			return
		}
		log.Printf("receive value:%v", message)
		data := []byte(time.Now().Format(common.DateTimeFormat))
		if _, err := conn.Write(data); err != nil {
			log.Println(err)
			return
		}
	}
}
