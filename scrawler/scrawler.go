package scrawler

import (
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"log/slog"

	"goScrawler/utils"

	"github.com/gocolly/colly/v2"
)

// Article structure for storing article data
type article struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Meta structure for storing metadata
type meta struct {
	Publisher string `json:"publisher"`
	PostAt    string `json:"post_at"`
	Views     int    `json:"views"`
}

// Content structure for storing complete page content
type content struct {
	URL         string  `json:"url"`
	ParentURL   string  `json:"parent_url,omitempty"`
	Meta        meta    `json:"meta"`
	Article     article `json:"article"`
	ScrapedAt   string  `json:"scraped_at"`
	ProcessedAt string  `json:"processed_at,omitempty"`
}

// Global shared variables
var (
	results        []content
	resultsMutex   = &sync.Mutex{}
	detailPagesMu  = &sync.Mutex{}
	catalogPagesMu = &sync.Mutex{}
	detailCount    = 0
	catalogCount   = 0
	successCount   = 0
	successCountMu = &sync.Mutex{}
)

// Start initiates the crawler with the given configuration
func Start(cfg *CrawlerConfig) []content {
	// Configure crawlers based on user input
	if cfg != nil {
		config.maxDetailPages = cfg.MaxDetailPages
		config.maxCatalogPages = cfg.MaxCatalogPages
	}

	// Set up parent collector for catalog pages
	parentCollector := setupParentCollector()
	// Set up child collector for detail pages
	childCollector := setupChildCollector()

	// Setup shared random user agent handling
	setupUserAgent(parentCollector)
	setupUserAgent(childCollector)

	// Setup error handling for both collectors
	setupErrorHandling(parentCollector, "Parent")
	setupErrorHandling(childCollector, "Child")

	// Configure parent collector to find and queue detail page URLs
	parentCollector.OnHTML(".news_title a", func(e *colly.HTMLElement) {
		catalogPagesMu.Lock()
		if catalogCount >= config.maxCatalogPages {
			catalogPagesMu.Unlock()
			return
		}
		catalogPagesMu.Unlock()

		detailPagesMu.Lock()
		if successCount >= config.maxDetailPages {
			detailPagesMu.Unlock()
			return
		}
		detailPagesMu.Unlock()

		detailURL := e.Request.AbsoluteURL(e.Attr("href"))
		parentURL := e.Request.URL.String()

		// Tell the child collector to visit this URL
		ctx := colly.NewContext()
		ctx.Put("parentURL", parentURL)
		_ = childCollector.Request("GET", detailURL, nil, ctx, nil)

		slog.Debug("Queued detail page",
			slog.String("url", detailURL),
			slog.String("parent_url", parentURL))
	})

	// Configure parent collector to find and follow next page links
	parentCollector.OnHTML(".page_nav a.next", func(e *colly.HTMLElement) {
		catalogPagesMu.Lock()
		if catalogCount >= config.maxCatalogPages {
			catalogPagesMu.Unlock()
			return
		}
		catalogCount++
		currentCatalogCount := catalogCount
		catalogPagesMu.Unlock()

		slog.Info("Processing catalog page",
			slog.Int("page_count", currentCatalogCount),
			slog.String("url", e.Request.URL.String()))

		nextURL := e.Request.AbsoluteURL(e.Attr("href"))
		if nextURL != "" {
			e.Request.Visit(nextURL)
			slog.Debug("Next page", slog.String("url", nextURL))
		} else {
			slog.Warn("No next page URL found")
		}
	})

	// Catalog page info extraction
	parentCollector.OnHTML("html", func(e *colly.HTMLElement) {
		recordCount := e.DOM.Find("em.all_count").Text()
		pageCount := e.DOM.Find("em.all_pages").Text()
		slog.Debug("Catalog page info",
			slog.String("record_count", recordCount),
			slog.String("page_count", pageCount),
			slog.String("url", e.Request.URL.String()))
	})

	// Configure child collector to process detail pages
	childCollector.OnHTML("html", func(e *colly.HTMLElement) {
		detailPagesMu.Lock()
		detailCount++
		currentDetailCount := detailCount
		detailPagesMu.Unlock()

		url := e.Request.URL.String()
		parentURL := e.Request.Ctx.Get("parentURL")

		slog.Debug("Processing detail page",
			slog.Int("page_count", currentDetailCount),
			slog.String("url", url),
			slog.String("parent_url", parentURL))

		// Extract content from detail page
		result, success := extractDetailPageContent(e, parentURL)

		if success {
			successCountMu.Lock()
			successCount++
			currentSuccessCount := successCount
			successCountMu.Unlock()

			resultsMutex.Lock()
			results = append(results, result)
			resultsMutex.Unlock()

			slog.Info("Successfully processed detail page",
				slog.Int("success_count", currentSuccessCount),
				slog.String("title", result.Article.Title),
				slog.String("url", url))
		} else {
			slog.Warn("Failed to extract content from detail page",
				slog.String("url", url),
				slog.String("parent_url", parentURL))
		}
	})

	// Start the crawling process
	slog.Info("Starting crawler", slog.String("base_url", config.baseURL))
	parentCollector.Visit(config.baseURL)

	// Wait for all collectors to finish
	parentCollector.Wait()
	childCollector.Wait()

	slog.Info("Crawling completed",
		slog.Int("catalog_pages_visited", catalogCount),
		slog.Int("detail_pages_visited", detailCount),
		slog.Int("successful_details", successCount),
		slog.Int("results_count", len(results)))

	return results
}

// setupUserAgent configures user agent rotation for a collector
func setupUserAgent(c *colly.Collector) {
	c.OnRequest(func(r *colly.Request) {
		randomUserAgent := config.userAgents[rand.Intn(len(config.userAgents))]
		r.Headers.Set("User-Agent", randomUserAgent)
		slog.Debug("Using User-Agent", slog.String("user_agent", randomUserAgent))
	})
}

// setupErrorHandling configures error handling for a collector
func setupErrorHandling(c *colly.Collector, collectorType string) {
	c.OnError(func(r *colly.Response, err error) {
		url := r.Request.URL.String()
		parentURL := r.Request.Ctx.Get("parentURL")

		slog.Error("Request failed",
			slog.String("collector", collectorType),
			slog.String("url", url),
			slog.String("parent_url", parentURL),
			slog.String("error", err.Error()))
	})
}

// extractDetailPageContent processes a detail page and extracts content
func extractDetailPageContent(e *colly.HTMLElement, parentURL string) (content, bool) {
	result := content{
		URL:         e.Request.URL.String(),
		ParentURL:   parentURL,
		ScrapedAt:   time.Now().Format("2006-01-02"),
		ProcessedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	// Extract article title
	title := e.ChildText("h1.arti_title")
	if title == "" {
		return result, false
	}
	result.Article.Title = strings.TrimSpace(title)

	// Extract article meta information
	result.Meta.PostAt = utils.ExtractDate(e.ChildText(".arti_metas .arti_update"))
	result.Meta.Publisher = strings.TrimSpace(e.ChildText(".arti_metas .arti_publisher"))

	viewsStr := e.ChildText(".arti_metas .arti_views > span")
	if viewsStr != "" {
		result.Meta.Views, _ = strconv.Atoi(viewsStr)
	}

	// Extract article body using colly's methods
	var bodyBuilder strings.Builder
	e.ForEach(".read p", func(_ int, paragraph *colly.HTMLElement) {
		paragraph.ForEach("span", func(_ int, span *colly.HTMLElement) {
			text := strings.TrimSpace(span.Text)
			if text != "" {
				bodyBuilder.WriteString(text)
				bodyBuilder.WriteString(" ")
			}
		})
		bodyBuilder.WriteString("\n")
	})

	result.Article.Body = bodyBuilder.String()

	// If we couldn't extract a body, check if we have at least a title
	return result, result.Article.Title != "" || result.Article.Body != ""
}
