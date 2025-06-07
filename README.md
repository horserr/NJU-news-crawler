# GoCrawler

A Parent-Child web crawler built in Go using the Colly framework. This crawler is designed to efficiently scrape websites by using separate crawlers for category pages and detail pages.

## Features

- Parent crawler traverses category/listing pages
- Child crawler extracts detailed information from detail pages
- Command-line arguments for customizing the crawl
- Different rate limiting for category vs. detail pages
- Robust error handling with comprehensive logging
- Results exported to JSON file

## Usage

Run the crawler with default settings:

```
go run main.go
```

### Command-line Arguments

- `-details` - Maximum number of detail pages to crawl (default: 40)
- `-catalogs` - Maximum number of catalog pages to crawl (default: 10)
- `-log` - Log level: debug, info, warn, error (default: info)
- `-output` - Output JSON file name (default: results.json)

Example:

```
go run main.go -details 100 -catalogs 5 -log debug -output results-2025.json
```

## Structure

- Parent crawler: Responsible for finding and following category pages and queuing detail page URLs
- Child crawler: Processes detail pages and extracts content

## Error Handling

Failed detail pages are logged with both their URLs and parent page URLs. Failed pages don't count toward the total number of crawled pages.

## Project Structure

- `main.go` - Entry point and command-line arguments
- `scrawler/`
  - `scrawler.go` - Main crawler implementation
  - `config.go` - Crawler configuration
- `utils/` - Utility functions

## Output

Results are saved to a JSON file with the following structure:

```json
[
  {
    "url": "https://example.com/page1",
    "parent_url": "https://example.com/list1",
    "meta": {
      "publisher": "Publisher Name",
      "post_at": "2025-05-20",
      "views": 123
    },
    "article": {
      "title": "Article Title",
      "body": "Article content..."
    },
    "scraped_at": "2025-05-20",
    "processed_at": "2025-05-20 15:04:05"
  },
  ...
]
```
