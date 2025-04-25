package main

import (
	"encoding/json"
	"fmt"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/gocolly/colly/v2"
)

func setupCollector() *colly.Collector {
	jar, _ := cookiejar.New(nil)

	// 创建 Colly Collector
	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(2),
	)
	c.SetCookieJar(jar)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 10,
		Delay:       300 * time.Millisecond,
		RandomDelay: 1 * time.Second,
	})

	return c
}

func main() {
	collector := setupCollector()
	startScraping(collector, BaseURL, MaxLinks)
	collector.Wait()

	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		panic(err)
	}
	// 将 JSON 数据写入文件
	fileName := "results.json"
	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JSON 数据已写入文件: %s\n", fileName)
}
