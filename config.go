package main

import (
	"goScrawler/utils"
	"os"
)

func init() {
	// 在 main 函数之前执行
	os.Setenv("GODEBUG", "tlsrsakex=1")
	os.Setenv("MY_DEBUG", "true")
	os.Setenv("MIN_LOG_LEVEL", "debug")

	utils.DebugLog(os.Stdout, utils.DEBUG, "Environment variable GODEBUG set")
}
