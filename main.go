package main

import (
	"encoding/json"
	"flag"
	"goScrawler/scrawler"
	"log/slog"
	"os"
)

func init() {
	// Set environment variable to handle older TLS implementations
	os.Setenv("GODEBUG", "tlsrsakex=1")
	slog.Debug("Environment variable GODEBUG set")
}

func main() {
	// Command-line flags
	maxDetailPages := flag.Int("details", 40, "Maximum number of detail pages to crawl")
	maxCatalogPages := flag.Int("catalogs", 10, "Maximum number of catalog pages to crawl")
	logLevel := flag.String("log", "info", "Log level (debug, info, warn, error)")
	outputFile := flag.String("output", "results.json", "Output file name")
	flag.Parse()

	// Configure logging
	var level slog.Level
	switch *logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// Start crawler with the provided configuration
	crawlerConfig := &scrawler.CrawlerConfig{
		MaxCatalogPages: *maxCatalogPages,
		MaxDetailPages:  *maxDetailPages,
	}

	slog.Info("Starting crawler",
		slog.Int("maxDetailPages", *maxDetailPages),
		slog.Int("maxCatalogPages", *maxCatalogPages))

	results := scrawler.Start(crawlerConfig)

	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal JSON data", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Write JSON data to file
	err = os.WriteFile(*outputFile, jsonData, 0644)
	if err != nil {
		slog.Error("Failed to write JSON data to file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("Write JSON data to file",
		slog.String("file_name", *outputFile),
		slog.Int("results_count", len(results)))
}
