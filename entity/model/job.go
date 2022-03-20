package model

import (
	"game/common"
	"time"
)

// Job
// cron表达式：分[0-59] 时[0-23] 日[1-31] 月[1-12] 周[0-6][SUN-SAT]
// 预定义规则：@hourly(每小时开始时)、@daily(每日开始时)、@weekly(每周开始时)、@monthly(每月开始时)、@yearly(每年开始时)
// 循环执行：@every Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
type Job struct {
	common.BaseUpdateModel
	Name        string    // 名称
	Code        string    // 编码
	Expression  string    // 表达式
	RunWay      uint8     // 执行方式：0.并行执行；1.串行执行；2.不执行。
	Method      string    // Timer类的执行方法
	State       uint8     // 状态：0.停用；1.启用。
	RunState    uint8     // 运行状态：0.未运行；1.运行中。
	EntryId     int       // 实例ID
	LastRunTime time.Time // 最后执行时间
}

func (t Job) TableName() string {
	return "sys_job"
}
