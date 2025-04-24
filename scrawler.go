// filepath: i:\Go\scrawler\scrawler.go
package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
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
	// 自定义处理子 URL 的逻辑
	title := e.DOM.Find("h1").Text()
	fmt.Printf("Title: %v\n", title)
	e.ForEach(".read p", func(_ int, para *colly.HTMLElement) {
		contentList := para.DOM.Find("span")

		if contentList.Size() > 0 {
			contentList.Each(func(_ int, span *goquery.Selection) {
				fmt.Printf("%v", span.Text())
			})
			fmt.Println()
		}
	})
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
	// c.Wait()
}
