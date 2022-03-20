package timer

import (
	"game/entity/model"
	"game/service/system"
	"game/utils/logutil"
	"github.com/robfig/cron/v3"
	"reflect"
)

var Timer *timer

var log = logutil.Log.WithField("kind", "timer")

// Start 定时任务入口
// https://darjun.github.io/2020/06/25/godailylib/cron/
func Start() {
	log.Info("定时任务启动_开始...")
	t := newTimer()
	// 查询配置的定时任务并启动
	jobs := system.GetJobList(model.Job{State: 1})
	if len(jobs) > 0 {
		v := reflect.ValueOf(t).Elem()
		for _, job := range jobs {
			m := v.MethodByName(job.Method)
			mv := m.Call(nil)[0]
			switch method := mv.Interface().(type) {
			case func():
				t.AddFunc(&job, method)
			case cron.Job:
				t.AddJob(&job, method)
			}
		}
	} else {
		log.Info("定时任务启动_未配置初始任务。")
	}
	// 启动定时任务
	t.start()
	log.Info("定时任务启动_完成！")
	Timer = t
}

type timer struct {
	cron      *cron.Cron // 按时执行，可并行执行
	cronDelay *cron.Cron // 延迟本次执行，直到上一次执行完成
	cronSkip  *cron.Cron // 若上次任务未执行完成，则跳过本次执行
}

func newTimer() *timer {
	logger := cron.VerbosePrintfLogger(log)
	logOp := cron.WithLogger(logger)
	recoverOp := cron.WithChain(cron.Recover(logger))
	delayOp := cron.WithChain(cron.DelayIfStillRunning(logger))
	skipOp := cron.WithChain(cron.SkipIfStillRunning(logger))
	return &timer{
		cron:      cron.New(logOp, recoverOp),
		cronDelay: cron.New(logOp, recoverOp, delayOp),
		cronSkip:  cron.New(logOp, recoverOp, skipOp),
	}
}

func (t *timer) start() {
	t.cron.Start()
	t.cronDelay.Start()
	t.cronSkip.Start()
}

func (t *timer) getCron(runWay uint8) *cron.Cron {
	switch runWay {
	case 1:
		return t.cronDelay
	case 2:
		return t.cronSkip
	default:
		return t.cron
	}
}

func (t *timer) AddFunc(job *model.Job, run func()) {
	eid, err := t.getCron(job.RunWay).AddFunc(job.Expression, run)
	t.addCallBack(job, eid, err)
}

func (t *timer) AddJob(job *model.Job, run cron.Job) {
	eid, err := t.getCron(job.RunWay).AddJob(job.Expression, run)
	t.addCallBack(job, eid, err)
}

func (t *timer) addCallBack(job *model.Job, eid cron.EntryID, err error) {
	if err != nil {
		log.Info(err)
		job.RunState = 0
	} else {
		job.EntryId = int(eid)
		job.RunState = 1
	}
	system.SaveJobById(job)
}

func (t *timer) Remove(jobId uint64) {
	job := system.GetJobById(jobId)
	t.getCron(job.RunWay).Remove(cron.EntryID(job.EntryId))
}
