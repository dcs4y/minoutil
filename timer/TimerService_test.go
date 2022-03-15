package timer

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	{
		// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
		d, err := time.ParseDuration("2s")
		if err != nil {
			log.Error(err)
		}
		// 延迟执行一次
		time.AfterFunc(d, func() {
			log.Info("============================")
		})
		// 自定义执行次数
		d, err = time.ParseDuration("3s")
		if err != nil {
			log.Error(err)
		}
		tt := time.NewTicker(d)
		for i := 0; i < 3; i++ {
			select {
			case tm := <-tt.C:
				log.Info(tm.String())
			}
		}
		tt.Stop()
	}
	{
		// 固定频率
		// 自定义日志
		logger := cron.VerbosePrintfLogger(log)
		c := cron.New(cron.WithLogger(logger))
		// cron表达式：分[0-59] 时[0-23] 日[1-31] 月[1-12] 周[0-6][SUN-SAT]
		c.AddFunc("30 3-6,20-23 * * *", func() {
			log.Info("On the half hour of 3-6am, 8-11pm")
		})
		// 预定义规则：@hourly(每小时开始时)、@daily(每日开始时)、@weekly(每周开始时)、@monthly(每月开始时)、@yearly(每年开始时)
		c.AddFunc("@hourly", func() {
			log.Info("Every hour")
		})
		// 循环执行：Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
		c.AddFunc("@every 1s", func() {
			time.Sleep(time.Second * 5)
			log.Info("tick every 1 second")
		})
		// Job接口。可给携带数据。
		{
			eid, err := c.AddJob("@every 1s", GreetingJob{"dj"})
			if err != nil {
				log.Info(err)
			}
			log.Info(eid)
			// 移除任务
			c.Remove(eid)
		}
		// 内置JobWrapper
		{
			// 捕获内部Job产生的 panic
			c.AddJob("@every 1s", cron.NewChain(cron.Recover(logger)).Then(&panicJob{}))
			// 触发时，如果上一次任务还未执行完成（耗时太长），则等待上一次任务完成之后再执行
			c.AddJob("@every 1s", cron.NewChain(cron.DelayIfStillRunning(logger)).Then(&delayJob{}))
			// 触发时，如果上一次任务还未完成，则跳过此次执行
			c.AddJob("@every 1s", cron.NewChain(cron.SkipIfStillRunning(logger)).Then(&skipJob{}))
		}
		c.Start()
	}
	{
		// 固定频率。cron表达式支持到秒。
		c := cron.New(cron.WithSeconds())
		c.Start()
	}
	time.Sleep(time.Second * 10)
}

func TestStart(t *testing.T) {
	Start()
	select {}
}

func Test_reflect(t *testing.T) {
	// 指针需要调用.Elem()
	v := reflect.ValueOf(&GreetingJob{Name: "dcs"}).Elem()
	fmt.Println(v)
	fmt.Println(v.Type())
	fmt.Println(v.NumField())
	m := v.MethodByName("Run")
	m.Call(nil)
}

type GreetingJob struct {
	Name string
}

func (g GreetingJob) Run() {
	log.Info("Hello ", g.Name)
}

type panicJob struct {
	count int
}

func (p *panicJob) Run() {
	p.count++
	if p.count == 1 {
		panic("oooooooooooooops!!!")
	}
	log.Info("hello world")
}

type delayJob struct {
	count int
}

func (d *delayJob) Run() {
	time.Sleep(2 * time.Second)
	d.count++
	log.Printf("%d: hello world\n", d.count)
}

type skipJob struct {
	count int32
}

func (d *skipJob) Run() {
	atomic.AddInt32(&d.count, 1)
	log.Printf("%d: hello world\n", d.count)
	if atomic.LoadInt32(&d.count) == 1 {
		time.Sleep(2 * time.Second)
	}
}
