package logutil

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

var log = logrus.New()

type LogConfig struct {
	LogLevel   string //日志级别：PANIC|FATAL|ERROR|WARN|INFO|DEBUG|TRACE
	FilePath   string //日志文件存储路径
	ConsoleLog bool   //是否打印控制台日志
}

func InitLog(config LogConfig) {
	// 设置日志级别
	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatalf("创建文件log.txt失败: %v", err)
	} else {
		logrus.SetLevel(logLevel) // 设置基础日志级别
		log.SetLevel(logLevel)    // 设置当前实例日志级别
	}
	// 在输出日志中添加文件名和方法信息
	//Log.SetReportCaller(true)
	// 设置日志输出格式
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.DateTime, // 设置json里的日期输出格式
	})
	// 重定向输出
	fileWrite := &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    20, // megabytes 兆M
		MaxBackups: 7,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}
	if config.ConsoleLog {
		log.SetOutput(io.MultiWriter(os.Stdout, fileWrite))
	} else {
		log.SetOutput(fileWrite)
	}
}

func GetLog(kind string) *logrus.Entry {
	return log.WithField("kind", kind)
}
