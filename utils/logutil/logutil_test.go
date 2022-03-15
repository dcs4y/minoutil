package logutil

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestLogger(t *testing.T) {
	// 定义默认key
	Log.WithFields(logrus.Fields{
		"默认Key": "内容",
		"age":   18,
	})
	Log.Info("初始化日志组件...")
	Log.Info("调用日志组件：")
	// 添加临时key
	Log.WithFields(logrus.Fields{
		"test": "content",
	}).Info("必填")
	// 调用没有默认key的日志对象记录日志
	Log.Info("同一个日志文件")
}
