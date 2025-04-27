package scrawler

import (
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
	"goScrawler/utils"
)

var (
	catalogPageCount = 0
	catalogMutex     = &sync.Mutex{}
)

func isCategoryPage(e *colly.HTMLElement) bool {
	urlStr := e.Request.URL.String()

	lastSegment, err := utils.GetLastSegment(urlStr)
	if err != nil || lastSegment == "" {
		return false
	}
	return strings.HasPrefix(lastSegment, "list")
}

func handleCategoryPage(e *colly.HTMLElement) {
	catalogMutex.Lock()
	if catalogPageCount >= config.maxCatalogPages {
		catalogMutex.Unlock()
		return
	}
	catalogPageCount++
	utils.DebugLog(os.Stdout, utils.DEBUG, "Processing catalog page: %vth %s\n", catalogPageCount, e.Request.URL.String())
	catalogMutex.Unlock()

	record_total_count := e.DOM.Find("em.all_count").Text()
	utils.DebugLog(os.Stdout, utils.INFO, "record_count: %v\n", record_total_count)

	page_count := e.DOM.Find("em.all_pages").Text()
	utils.DebugLog(os.Stdout, utils.INFO, "page_count: %v\n", page_count)

	e.ForEach(".news_title a", func(_ int, el *colly.HTMLElement) {
		e.Request.Visit(e.Request.AbsoluteURL(el.Attr("href")))
	})

	nextPage := e.DOM.Find(".page_nav a.next").AttrOr("href", "")
	if nextPage == "" {
		panic("Next page")
	}
	e.Request.Visit(e.Request.AbsoluteURL(nextPage))
	utils.DebugLog(os.Stdout, utils.DEBUG, "Next page: %v\n", nextPage)
}
