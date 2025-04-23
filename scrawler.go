package main

import (
    "fmt"
    "os"

    "github.com/gocolly/colly/v2"
)

func setupCollector() *colly.Collector {
    // 设置 GODEBUG 环境变量
    os.Setenv("GODEBUG", "tlsrsakex=1")

    // 创建 Colly Collector
    c := colly.NewCollector(
        colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
    )

    // 设置 HTML 处理逻辑
    c.OnHTML("h1", func(e *colly.HTMLElement) {
        content := e.Text
        fmt.Printf("content: %v\n", content)
    })

    // 设置错误处理逻辑
    c.OnError(func(r *colly.Response, err error) {
        fmt.Fprintf(os.Stderr, "Request URL: %v failed with response: %v, error: %v\n", r.Request.URL, r, err)
    })

    return c
}

func startScraping(c *colly.Collector, url string) {
    c.Visit(url)
}

func main() {
    collector := setupCollector()
    startScraping(collector, "https://jw.nju.edu.cn/88/0f/c26263a755727/page.htm")
}