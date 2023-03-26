package collyutil

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

type CollectorConfig struct {
	AllowUrlRevisit      bool                                  // 是否允许相同url重复请求
	Async                bool                                  // 是否异步抓取
	LimitRuleDomainGlob  string                                // glob模式匹配域名
	LimitRuleParallelism int                                   // 匹配到的域名的并发请求数
	LimitRuleRandomDelay int                                   // 在发起一个新请求时的随机等待时间(秒)
	Proxy                func(*http.Request) (*url.URL, error) // 代理方法
}

func NewCollectorConfig() *CollectorConfig {
	cc := &CollectorConfig{
		AllowUrlRevisit:      false,
		Async:                true,
		LimitRuleDomainGlob:  "",
		LimitRuleParallelism: 10,
		LimitRuleRandomDelay: 10,
		Proxy:                http.ProxyFromEnvironment,
	}
	return cc
}

// NewCollector 生成一个collector对象
func NewCollector(collectorConfig *CollectorConfig) *colly.Collector {
	if collectorConfig == nil {
		collectorConfig = NewCollectorConfig()
	}
	collector := colly.NewCollector()
	collector.WithTransport(&http.Transport{
		Proxy: collectorConfig.Proxy,
		DialContext: (&net.Dialer{
			Timeout:   300 * time.Second,
			KeepAlive: 300 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       300 * time.Second,
		TLSHandshakeTimeout:   300 * time.Second,
		ExpectContinueTimeout: 300 * time.Second,
	})

	//是否允许相同url重复请求
	collector.AllowURLRevisit = collectorConfig.AllowUrlRevisit

	//默认是同步,配置为异步,这样会提高抓取效率
	collector.Async = collectorConfig.Async

	collector.DetectCharset = true

	// 对于匹配的域名(当前配置为任何域名),将请求并发数配置为2
	// 通过测试发现,RandomDelay参数对于同步模式也生效
	if err := collector.Limit(&colly.LimitRule{
		// glob模式匹配域名
		DomainGlob: collectorConfig.LimitRuleDomainGlob,
		// 匹配到的域名的并发请求数
		Parallelism: collectorConfig.LimitRuleParallelism,
		// 在发起一个新请求时的随机等待时间
		RandomDelay: time.Duration(collectorConfig.LimitRuleRandomDelay) * time.Second,
	}); err != nil {
		log.Println("生成一个collector对象, 限速配置失败", err)
	}

	//配置反爬策略(设置ua和refer扩展)
	extensions.RandomUserAgent(collector)
	extensions.Referer(collector)
	return collector
}
