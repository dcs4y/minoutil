package cacheclient

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"testing"
	"time"
)

func Test_cache(t *testing.T) {
	//go-cache是内存中的键：类似于memcached的值存储/缓存，适用于在一台计算机上运行的应用程序。它的主要优点是，
	//由于本质上是map[string]interface{}具有到期时间的线程安全的，因此不需要通过网络序列化或传输其内容。
	//可以在给定的持续时间内或永久存储任何对象，并且可以由多个goroutine安全使用缓存。
	//尽管不打算将go-cache用作持久性数据存储，
	//但可以将整个缓存保存到文件中并从文件中加载（c.Items()用于检索要映射的项目映射并NewFrom()从反序列化的缓存中创建缓存）以进行恢复从停机时间很快。

	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cache.New(5*time.Minute, 10*time.Minute)

	// Set the value of the key "foo" to "bar", with the default expiration time
	c.Set("foo", "bar", cache.DefaultExpiration)

	// Set the value of the key "baz" to 42, with no expiration time
	// (the item won't be removed until it is re-set, or removed using
	// c.Delete("baz")
	c.Set("baz", 42, cache.NoExpiration)

	// Get the string associated with the key "foo" from the cache
	{
		foo, found := c.Get("foo")
		if found {
			fmt.Println(foo)
		}
	}
}

func Test_util(t *testing.T) {
	Set("k1", "v1")
	SetAlways("k2", 222)
	fmt.Println(Get("k2"))
	Delete("k2")
	fmt.Println(Get("k2"))
	fmt.Println(Get("k1"))
	<-time.After(5 * time.Minute)
	fmt.Println(Get("k1"))
}
