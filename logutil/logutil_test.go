package logutil

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestLogger(t *testing.T) {
	log := GetLog("test")
	// 定义默认key
	log.WithFields(logrus.Fields{
		"默认Key": "内容",
		"age":   18,
	})
	log.Info("初始化日志组件...")
	log.Info("调用日志组件：")
	// 添加临时key
	log.WithFields(logrus.Fields{
		"test": "content",
	}).Info("必填")
	// 调用没有默认key的日志对象记录日志
	log.Info("同一个日志文件")
}
