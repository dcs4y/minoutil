package limiterutil

import (
	"golang.org/x/time/rate"
	"sync"
)

var limiterMap = make(map[string]*rate.Limiter)

var lock sync.RWMutex

// NewLimiter 限流器
func NewLimiter(name string, r rate.Limit, b int) *rate.Limiter {
	lock.Lock()
	defer lock.Unlock()
	if limiter, ok := limiterMap[name]; ok {
		return limiter
	}
	//参数r Limit。代表每秒可以向Token桶中产生多少token。Limit实际上是float64的别名。
	//参数b int。b代表Token桶的容量大小。 那么，对于NewLimiter(10, 1)来说，其构造出的限流器含义为，其令牌桶大小为1, 以每秒10个Token的速率向桶中放置Token。
	limiter := rate.NewLimiter(r, b)
	limiterMap[name] = limiter
	limiter.Limit()
	return limiter
}

func GetLimiter(name string) *rate.Limiter {
	lock.Lock()
	defer lock.Unlock()
	return limiterMap[name]
}
