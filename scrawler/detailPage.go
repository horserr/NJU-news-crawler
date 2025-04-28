package scrawler

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"goScrawler/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
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
	meta := buildMeta(e)
	article := buildArticle(e)
	content := content{
		URL:       e.Request.URL.String(),
		Meta:      meta,
		Article:   article,
		ScrapedAt: time.Now().Format("2006-01-02"),
	}

	resultsMutex.Lock()
	results = append(results, content)
	resultsMutex.Unlock()
}

func buildMeta(e *colly.HTMLElement) meta {
	meta := meta{}
	meta.PostAt = utils.ExtractDate(e.DOM.Find(".arti_metas .arti_update").Text())
	meta.Publisher = e.DOM.Find(".arti_metas .arti_publisher").Text()
	meta.Views, _ = strconv.Atoi(e.DOM.Find(".arti_metas .arti_views > span").Text())
	return meta
}

func buildArticle(e *colly.HTMLElement) article {
	article := article{}
	article.Title = e.DOM.Find("h1.arti_title").Text()

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
	article.Body = textContent

	return article
}
