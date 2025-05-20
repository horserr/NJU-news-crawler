package scrawler

import (
	"net/http/cookiejar"
	"time"

	"github.com/gocolly/colly/v2"
)

// Config 结构体定义
type Config struct {
	userAgents      []string
	maxCatalogPages int
	maxDetailPages  int
	maxDepth        int
	baseURL         string
}

// 初始化配置
var config = Config{
	userAgents: []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPad; CPU OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
	},
	maxCatalogPages: 10,
	maxDetailPages:  40,
	baseURL:         "https://jw.nju.edu.cn/ggtz/list.htm",
}

func setupCollector() *colly.Collector {
	jar, _ := cookiejar.New(nil)

	// 创建 Colly Collector
	c := colly.NewCollector(
		colly.Async(true),
		// colly.MaxDepth(10),
	)
	c.SetCookieJar(jar)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 12,
		Delay:       300 * time.Millisecond,
		RandomDelay: 250 * time.Millisecond,
	})

	return c
}
