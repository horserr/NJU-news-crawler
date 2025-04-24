package main

import (
	"net/http/cookiejar"
	"time"

	"github.com/gocolly/colly/v2"
)

func setupCollector() *colly.Collector {
	jar, _ := cookiejar.New(nil)

	// 创建 Colly Collector
	c := colly.NewCollector()
	c.SetCookieJar(jar)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 10,
		Delay:       300 * time.Millisecond,
		RandomDelay: 1 * time.Second,
	})

	return c
}

func main() {
	collector := setupCollector()
	startScraping(collector, BaseURL, MaxLinks)
}
