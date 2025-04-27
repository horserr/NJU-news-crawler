package main

import (
	"encoding/json"
	"goScrawler/scrawler"
	"goScrawler/utils"
	"os"
)

func main() {
	results := scrawler.Start()

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
