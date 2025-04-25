// filepath: i:\Go\scrawler\scrawler.go
package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type Content struct {
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

var (
	results []Content  // 存储所有结果
	mutex   sync.Mutex // 用于保护 results 的并发写入
)

func handleBaseHTML(e *colly.HTMLElement, baseURL string, maxLinks int) {
	e.ForEach(".news_title a", func(i int, el *colly.HTMLElement) {
		if i >= maxLinks {
			return // 达到最大链接数，停止处理
		}

		// 获取相对链接
		relativeLink := el.Attr("href")

		// 将相对链接转换为绝对链接
		base, _ := url.Parse(baseURL)
		absoluteLink := base.ResolveReference(&url.URL{Path: relativeLink}).String()
		debugLog(os.Stderr, "absoluteLink: %v\n", absoluteLink)

		// 访问新的链接
		el.Request.Visit(absoluteLink)
	})
}
func handleSubHTML(e *colly.HTMLElement) {
	// 获取当前页面的 URL
	currentURL := e.Request.URL.String()

	// 获取创建者信息（假设从 meta 标签中获取）
	created_at := e.DOM.Find(".arti_metas .arti_update").Text()
	created_at = extractDate(created_at) // 提取日期

	debugLog(os.Stderr, "created_at: %v\n", created_at)

	// 获取页面内容
	var textContent string
	e.ForEach(".read p", func(_ int, para *colly.HTMLElement) {
		contentList := para.DOM.Find("span")

		if contentList.Size() > 0 {
			contentList.Each(func(_ int, span *goquery.Selection) {
				textContent += span.Text() + " "
			})
			textContent += "\n"
		}
	})
	// 创建结果对象
	result := Content{
		URL:       currentURL,
		CreatedAt: created_at,
		Text:      textContent,
	}

	// 使用 Mutex 保护并发写入
	mutex.Lock()
	results = append(results, result)
	mutex.Unlock()
}

func startScraping(c *colly.Collector, baseURL string, maxLinks int) {
	c.OnRequest(func(r *colly.Request) {
		// 随机选择一个 User-Agent
		randomUserAgent := UserAgents[rand.Intn(len(UserAgents))]
		r.Headers.Set("User-Agent", randomUserAgent)
		// fmt.Printf("Using User-Agent: %s\n", randomUserAgent)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Fprintf(os.Stderr, "Request URL: %v failed with response: %v, error: %v\n", r.Request.URL, r, err)
	})

	// 注册一次 OnHTML，根据深度调用不同的处理逻辑
	c.OnHTML("html", func(e *colly.HTMLElement) {
		if e.Request.Depth == 1 {
			// Base URL 的处理逻辑
			handleBaseHTML(e, baseURL, maxLinks)
		} else {
			// 子 URL 的处理逻辑
			handleSubHTML(e)
		}
	})

	c.Visit(baseURL)
	c.Wait()
}
