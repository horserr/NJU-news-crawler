// This file is deprecated. Functionality has been moved to scrawler.go
package scrawler

import (
	"strings"

	"goScrawler/utils"

	"github.com/gocolly/colly/v2"
)

// Legacy code - kept for reference only
// Use the parent-child pattern in scrawler.go instead

func isDetailPage(e *colly.HTMLElement) bool {
	urlStr := e.Request.URL.String()

	lastSegment, err := utils.GetLastSegment(urlStr)
	if err != nil || lastSegment == "" {
		return false
	}
	return strings.HasPrefix(lastSegment, "page")
}
