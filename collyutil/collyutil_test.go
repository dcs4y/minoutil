package collyutil

import (
	"github.com/gocolly/colly"
	"go.uber.org/zap"
	"log"
	"strings"
	"testing"
)

func Test_colly(t *testing.T) {
	var ss = &spiderService{}
	err := ss.StartCrawNovelListTask()
	if err != nil {
		t.Log(err)
	}
}

type spiderService struct {
	spiderConfig         *CollectorConfig
	novelListCollector   *colly.Collector
	chapterListCollector *colly.Collector
	chapterCollector     *colly.Collector
}
type novel struct {
	Title                    string
	Author                   string
	Category                 string
	Summary                  string
	ChapterCount             int
	WordCount                string
	CoverSrcUrl              string
	NovelSrcUrl              string
	CurrentCrawChapterPageNo int
}
type chapter struct {
	Novel         *novel
	Title         string
	ChapterSrcUrl string
	Content       string
	Sort          int
}

/*
*
初始化collector
*/
func (ss *spiderService) initCollector() {
	ss.configNovelListCollector()
	ss.configChapterListCollector()
	ss.configChapterCollector()
}

func (ss *spiderService) NewCollector() *colly.Collector {
	ss.spiderConfig = NewCollectorConfig()
	collector := NewCollector(ss.spiderConfig)
	// 配置IP代理
	//collector.SetProxyFunc(NewProxyFunc(""))
	return collector
}

/*
*
配置NovelListCollector
*/
func (ss *spiderService) configNovelListCollector() {
	//避免对collector对象的每个回调注册多次, 否则回调内逻辑重复执行多次, 会引发逻辑错误
	if ss.novelListCollector != nil {
		return
	}
	ss.novelListCollector = ss.NewCollector()

	ss.novelListCollector.OnHTML("div.list_main li", func(element *colly.HTMLElement) {
		// 抽取某小说的入口页面地址和章节列表页的入口地址
		novelUrl, exist := element.DOM.Find("div.book-img-box a").Attr("href")
		if !exist {
			log.Println("爬取小说列表页, 抽取当前小说的入口url, 异常", zap.Any("novelUrl", novelUrl))
			return
		}
		chapterListUrl := strings.ReplaceAll(novelUrl, "book", "chapter")
		log.Println("爬取小说列表页, 抽取章节列表的入口url, 完成", zap.Any("chapterListUrl", chapterListUrl))

		//抽取小说剩余信息，并组装novel对象
		novel := &novel{}
		novel.Title = strings.TrimSpace(element.DOM.Find("div.book-mid-info p.t").Text())
		novel.NovelSrcUrl = chapterListUrl
		novel.CoverSrcUrl = element.DOM.Find("div.book-img-box img").AttrOr("src", "")
		novel.Author = strings.TrimSpace(element.DOM.Find("div.book-mid-info p.author span").First().Text())
		novel.Category = strings.TrimSpace(element.DOM.Find("div.book-mid-info p.author a").Text())
		novel.Summary = strings.TrimSpace(element.DOM.Find("div.book-mid-info p.intro").Text())
		novel.WordCount = strings.TrimSpace(element.DOM.Find("div.book-mid-info p.update").Text())

		// 创建上下文对象
		ctx := colly.NewContext()
		ctx.Put("novel", novel)

		// 爬取章节列表页
		log.Println("爬取小说列表页, 开始", zap.Any("novelTitle", novel.Title), zap.Any("chapterListUrl", chapterListUrl))
		if err := ss.chapterListCollector.Request("GET", chapterListUrl, nil, ctx, nil); err != nil {
			log.Println("爬取小说列表页, 爬取章节列表页, 异常", zap.Any("chapterListUrl", chapterListUrl))
			return
		}
	})

	/**
	爬取当前列表页的下一页
	*/
	ss.novelListCollector.OnHTML("div.tspage a.next", func(element *colly.HTMLElement) {
		nextUrl := element.Request.AbsoluteURL(element.Attr("href"))
		log.Println("爬取小说列表页的下一页, 开始", zap.Any("nextUrl", nextUrl))

		if err := ss.novelListCollector.Visit(nextUrl); err != nil {
			log.Println("爬取小说列表页的下一页, 异常", zap.Any("nextUrl", nextUrl), zap.Error(err))
			return
		}

		log.Println("爬取小说列表页的下一页, 完成", zap.Any("nextUrl", nextUrl))
	})

	ss.novelListCollector.OnError(func(response *colly.Response, e error) {
		log.Println("爬取小说列表页, OnError", zap.Any("url", response.Request.URL.String()), zap.Error(e))

		//请求重试
		response.Request.Retry()
	})

	log.Println("配置NovelListCollector, 完成")
}

/*
*
配置ChapterListCollector
*/
func (ss *spiderService) configChapterListCollector() {
	if ss.chapterListCollector != nil {
		return
	}
	ss.chapterListCollector = ss.NewCollector()

	ss.chapterListCollector.OnRequest(func(r *colly.Request) {
		log.Println("爬取章节列表页, OnRequest", zap.Any("url", r.URL.String()))
	})
	// 从章节列表页抓取第一章节的入口地址
	ss.chapterListCollector.OnHTML("div.catalog_b li:nth-child(1) a", func(h *colly.HTMLElement) {
		// 抽取某章节的地址
		chapterUrl, exist := h.DOM.Attr("href")
		if !exist {
			log.Println("爬取章节列表页, 爬取第1章, 抽取chapterUrl, 异常", zap.Any("srcUrl", h.Request.URL))
			return
		}
		chapterUrl = h.Request.AbsoluteURL(chapterUrl)
		chapterTitle := h.DOM.Text()
		log.Println("爬取章节列表页, 爬取第1章, 抽取chapterUrl, 完成", zap.Any("chapterUrl", chapterUrl), zap.Any("chapterTitle", chapterTitle))

		// 获取上下文信息
		novel := h.Response.Ctx.GetAny("novel").(*novel)
		novel.ChapterCount = h.DOM.Parent().Parent().Find("li").Length()
		novel.CurrentCrawChapterPageNo = 0

		// 爬取章节
		log.Println("爬取章节列表页, 开始爬取第1章", zap.Any("novelTitle", novel.Title), zap.Any("chapterTitle", chapterTitle))
		if err := ss.chapterCollector.Request("GET", chapterUrl, nil, h.Response.Ctx, nil); err != nil {
			log.Println("爬取章节列表页, 爬取第1章, 异常", zap.Any("chapterUrl", chapterUrl), zap.Error(err))
			return
		}
	})
	ss.chapterListCollector.OnError(func(response *colly.Response, e error) {
		log.Println("爬取章节列表页, OnError", zap.Any("url", response.Request.URL.String()), zap.Error(e))

		//请求重试
		response.Request.Retry()
	})
}

/*
*
配置configChapterCollector
*/
func (ss *spiderService) configChapterCollector() {
	if ss.chapterCollector != nil {
		return
	}
	ss.chapterCollector = ss.NewCollector()

	// 爬取章节
	ss.chapterCollector.OnHTML("div.mlfy_main", func(h *colly.HTMLElement) {
		chapterTitle := strings.TrimSpace(h.DOM.Find("h3.zhangj").Text())
		content, err := h.DOM.Find("div.read-content").Html()
		if err != nil {
			log.Println("爬取章节, 解析内容, 异常", zap.Error(err))
			return
		}

		// 获取上下文信息
		novel := h.Response.Ctx.GetAny("novel").(*novel)
		// 累加爬取的章节页码
		novel.CurrentCrawChapterPageNo++

		chapter := &chapter{}
		chapter.Content = content
		chapter.Novel = novel
		chapter.Title = chapterTitle
		chapter.ChapterSrcUrl = h.Request.URL.String()
		chapter.Sort = novel.CurrentCrawChapterPageNo

		log.Println("爬取章节, 完成", zap.Any("novelTitle", chapter.Novel.Title), zap.Any("chapterTitle", chapter.Title), zap.Any("novelSrcUrl", chapter.Novel.NovelSrcUrl), zap.Any("chapterSrcUrl", chapter.ChapterSrcUrl), zap.Any("chapter", chapter))
	})
	//通过翻页按钮爬取下一章
	ss.chapterCollector.OnHTML("p.mlfy_page a:contains(下一章)", func(h *colly.HTMLElement) {
		nextChapterUrl, exist := h.DOM.Attr("href")
		if !exist {
			log.Println("爬取下一章, 抽取下一页url， 异常", zap.Any("currentPage", h.Request.URL.String()))
			return
		}

		log.Println("爬取下一章, 开始爬取", zap.Any("currentPage", h.Request.URL.String()), zap.Any("nextChapterUrl", nextChapterUrl))
		if err := ss.chapterCollector.Request("GET", nextChapterUrl, nil, h.Response.Ctx, nil); err != nil {
			log.Println("爬取下一章, 异常", zap.Any("currentPage", h.Request.URL.String()), zap.Any("nextChapterUrl", nextChapterUrl))
			return
		}
	})
	ss.chapterCollector.OnError(func(response *colly.Response, e error) {
		log.Println("爬取章节, OnError", zap.Any("url", response.Request.URL.String()), zap.Error(e))

		//请求重试
		response.Request.Retry()
	})
	ss.chapterCollector.OnResponse(func(r *colly.Response) {
		log.Println("爬取章节, OnResponse, 保存文件", zap.Any("url", r.Request.URL.String()))
	})
}

/*
*
启动小说列表页爬取任务
*/
func (ss *spiderService) StartCrawNovelListTask() error {
	// 初始化collector
	ss.initCollector()

	if err := ss.novelListCollector.Visit("https://www.517shu.com/sort_2"); err != nil {
		log.Println("启动小说列表页爬取任务, 异常", zap.Error(err))
		return err
	}

	//若开启异步爬取模式, 则等待爬取线程执行完成
	if ss.spiderConfig.Async {
		log.Println("启动小说列表页爬取任务, 等待线程执行完成")
		ss.novelListCollector.Wait()
	}

	log.Println("启动小说列表页爬取任务, 完成")
	return nil
}
