package scrawler

import (
	"goScrawler/utils"
	"math/rand"
	"os"
	"sync"

	"github.com/gocolly/colly/v2"
)

var (
	results      []Content
	resultsMutex = &sync.Mutex{}
)

type Content struct {
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

type pageHandler struct {
	isPageType func(*colly.HTMLElement) bool
	handler    func(*colly.HTMLElement)
}

func Start() []Content {
	c := setupCollector()

	// 使用切片存储页面类型判断函数和处理函数的映射
	handlers := []pageHandler{
		{isPageType: isCategoryPage, handler: handleCategoryPage},
		{isPageType: isDetailPage, handler: handleDetailPage},
	}

	c.OnRequest(func(r *colly.Request) {
		randomUserAgent := config.userAgents[rand.Intn(len(config.userAgents))]
		r.Headers.Set("User-Agent", randomUserAgent)
		utils.DebugLog(os.Stdout, utils.DEBUG, "Using User-Agent: %s", randomUserAgent)
	})

	c.OnError(func(r *colly.Response, err error) {
		utils.DebugLog(os.Stderr, utils.ERROR, "Request URL: %v failed with response: %v, error: %v\n", r.Request.URL, r, err)
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		for _, handler := range handlers {
			if handler.isPageType(e) {
				handler.handler(e)
				return
			}
		}
		utils.DebugLog(os.Stderr, utils.WARN, "No handler found for URL: %s\n", e.Request.URL.String())
	})

	c.Visit(config.baseURL)
	c.Wait()
	return results
}
