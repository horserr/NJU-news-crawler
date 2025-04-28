package main

import (
	"encoding/json"
	"goScrawler/scrawler"
	"goScrawler/utils"
	"os"
	"regexp"
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

	utils.DebugLog(os.Stdout, utils.INFO, "Write JSON data to file: %s", fileName)
}
