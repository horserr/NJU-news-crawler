package main

import (
	"log/slog"
	"os"
)

func init() {
	// 在 main 函数之前执行
	os.Setenv("GODEBUG", "tlsrsakex=1")

	slog.Debug("Environment variable GODEBUG set")
}
