package scrawler

import (
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"goScrawler/utils"
)

var (
	detailPageCount = 0
	detailMutex     = &sync.Mutex{}
)

func isDetailPage(e *colly.HTMLElement) bool {
	urlStr := e.Request.URL.String()

	lastSegment, err := utils.GetLastSegment(urlStr)
	if err != nil || lastSegment == "" {
		return false
	}
	return strings.HasPrefix(lastSegment, "page")
}

func handleDetailPage(e *colly.HTMLElement) {
	// 检查内容详情页计数器
	detailMutex.Lock()
	if detailPageCount >= config.maxDetailPages {
		detailMutex.Unlock()
		return
	}
	detailPageCount++
	detailMutex.Unlock()

	// 获取当前页面的 URL
	currentURL := e.Request.URL.String()
	created_at := e.DOM.Find(".arti_metas .arti_update").Text()
	created_at = utils.ExtractDate(created_at) // 提取日期

	utils.DebugLog(os.Stdout, utils.INFO, "Processed detail page: %s\n", currentURL)

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

	// 使用 Mutex 保护并发写入
	resultsMutex.Lock()
	results = append(results, Content{
		URL:       currentURL,
		CreatedAt: created_at,
		Text:      textContent,
	})
	resultsMutex.Unlock()
}
