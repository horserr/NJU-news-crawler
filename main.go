package main

import (
	"encoding/json"
	"goScrawler/scrawler"
	"log/slog"
	"os"
)

func main() {
	results := scrawler.Start()

	// re := regexp.MustCompile(`[\p{C}]`)
	// for i := range results {
	// 	results[i].Article.Title = re.ReplaceAllString(results[i].Article.Title, "")
	// }

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

	slog.Info("Write JSON data to file", slog.String("file_name", fileName))
}
