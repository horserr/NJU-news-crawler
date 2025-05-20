package scrawler

import (
	"log/slog"
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
	slog.Debug("Processing catalog page", slog.Int("page_count", catalogPageCount), slog.String("url", e.Request.URL.String()))
	catalogMutex.Unlock()

	record_total_count := e.DOM.Find("em.all_count").Text()
	slog.Debug("Record count", slog.String("count", record_total_count))

	page_count := e.DOM.Find("em.all_pages").Text()
	slog.Debug("Page count", slog.String("count", page_count))

	e.ForEach(".news_title a", func(_ int, el *colly.HTMLElement) {
		e.Request.Visit(e.Request.AbsoluteURL(el.Attr("href")))
	})

	nextPage := e.DOM.Find(".page_nav a.next").AttrOr("href", "")
	if nextPage == "" {
		panic("Next page")
	}
	e.Request.Visit(e.Request.AbsoluteURL(nextPage))
	slog.Debug("Next page", slog.String("url", nextPage))
}
