package timer

import (
	"fmt"
	"game/utils/dingdingclient"
	"github.com/robfig/cron/v3"
)

var (
	bootUrl = "https://oapi.dingtalk.com/robot/send?access_token=ca6441d14175831d2f1e4e5409421b7a1a7859824c2cbb59cd6a705f1f318928"
	secret  = "SEC8864fda8960e963bbac681d27aab862957666aa97b5f0e449e9082a6fd6b26ea"
	robot   = dingdingclient.NewWebhook(bootUrl, secret)
)

// ===========================以下为根据业务的执行任务，需要在数据库进行相应配置========================== //

// RunFuncTemp 方法形式示例
func (t timer) RunFuncTemp() func() {
	return func() {
		fmt.Println("Run1")
	}
}

// RunJobTemp 任务形式示例
func (t timer) RunJobTemp() cron.Job {
	return RunJob{Name: "dong"}
}

type RunJob struct {
	Name string
}

func (j RunJob) Run() {
	log.Info("Hello ", j.Name)
}

// ParseLianJiaJob 链家成交数据解析
func (t timer) ParseLianJiaJob() func() {
	return func() {
		/*robot.Send(nil, &dingdingclient.TextBody{Text: "抓取成交数据开始："})
		// 主站网址
		homeUrl := "https://cq.lianjia.com/"
		// 抓取类型
		typeUri := "chengjiao"
		// 月份，如：202212
		now := time.Now()
		now = now.AddDate(0, -1, 0)
		month := now.Format("200601")
		robot.Send(nil, &dingdingclient.TextBody{Text: "抓取成交数据结束！"})*/
	}
}
