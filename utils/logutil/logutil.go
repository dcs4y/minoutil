package logutil

import (
	"game/common"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

var Log = logrus.New()

func init() {
	// 设置日志级别
	//log.SetLevel(logrus.TraceLevel)
	// 在输出日志中添加文件名和方法信息
	//Log.SetReportCaller(true)
	// 设置日志输出格式
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: common.DateTimeFormat, // 设置json里的日期输出格式
	})
	// 重定向输出
	fileWrite, err := os.OpenFile(filepath.Join(common.DataPath, "log", "log.txt"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		Log.Fatalf("创建文件log.txt失败: %v", err)
	}
	if common.WS.ServerConfig.Active == common.ServerActiveDev {
		Log.SetOutput(io.MultiWriter(os.Stdout, fileWrite))
	} else {
		Log.SetOutput(fileWrite)
	}
}
