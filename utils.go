package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func debugLog(f *os.File, format string, args ...interface{}) {
	if DEBUG {
		fmt.Fprintf(f, "[DEBUG] "+format+"\n", args...)
	}
}

func extractDate(input string) string {
	// 定义正则表达式，匹配类似 "2025-04-24" 的日期格式
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	// 提取匹配的日期
	match := re.FindString(input)
	return strings.TrimSpace(match)
}
