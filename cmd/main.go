package main

import "github.com/LucasZatta/letter-goxd/internal/lists"

func main() {
	scraper := lists.New()
	_ = scraper.ScrapeWatchlist("zvttx")
}
