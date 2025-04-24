package main

import (
	"fmt"
	"os"
)

func debugLog(f *os.File, format string, args ...interface{}) {
	if debug {
		fmt.Fprintf(f, "[DEBUG] "+format+"\n", args...)
	}
}
