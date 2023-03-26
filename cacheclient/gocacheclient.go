package cacheclient

import (
	"github.com/patrickmn/go-cache"
	"time"
)

// Cache 带过期时间的简单缓存
var Cache *cache.Cache

func init() {
	Cache = cache.New(5*time.Minute, 10*time.Minute)
}

// Set 设置缓存值，默认过期时间。
func Set(k string, x interface{}) {
	Cache.SetDefault(k, x)
}

// SetWithTimeout 设置缓存值，并设置过期时间。
func SetWithTimeout(k string, x interface{}, d time.Duration) {
	Cache.Set(k, x, d)
}

// SetAlways 设置缓存值，永不过期。
func SetAlways(k string, x interface{}) {
	Cache.Set(k, x, cache.NoExpiration)
}

// Get 获取已有的缓存数据
func Get(k string) (interface{}, bool) {
	return Cache.Get(k)
}

// Delete 删除一个缓存数据
func Delete(k string) {
	Cache.Delete(k)
}
