package utils

import (
	"net/url"
	"regexp"
	"strings"
)

func ExtractDate(input string) string {
	// 定义正则表达式，匹配类似 "2025-04-24" 的日期格式
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	// 提取匹配的日期
	match := re.FindString(input)
	return strings.TrimSpace(match)
}

func GetLastSegment(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	pathSegments := strings.Split(parsedURL.Path, "/")
	if len(pathSegments) == 0 {
		return "", nil
	}
	lastSegment := pathSegments[len(pathSegments)-1]
	return lastSegment, nil
}
