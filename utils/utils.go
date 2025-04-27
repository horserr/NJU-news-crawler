package utils

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	default:
		return "UNKNOWN"
	}
}

var (
	_Debug      = getEnvAsBool("MY_DEBUG", false)
	minLogLevel = getEnvAsLogLevel("MIN_LOG_LEVEL", INFO)
)

func getEnvAsBool(name string, defaultValue bool) bool {
	valStr := os.Getenv(name)
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}

func getEnvAsLogLevel(name string, defaultValue LogLevel) LogLevel {
	valStr := strings.ToLower(os.Getenv(name)) // 转为小写以便匹配
	if valStr == "" {
		return defaultValue
	}
	switch valStr {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warning", "warn": // 支持 "warning" 和 "warn"
		return WARN
	case "error":
		return ERROR
	default:
		return defaultValue // 如果不匹配，返回默认值
	}
}

func DebugLog(w io.Writer, level LogLevel, format string, args ...any) {
	if !_Debug || level < minLogLevel {
		return
	}
	fmt.Fprintf(w, "[%s] %s\n", level.String(), fmt.Sprintf(format, args...))
}

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