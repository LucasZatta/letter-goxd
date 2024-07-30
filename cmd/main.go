package main

import (
	"fmt"
	"time"

	"github.com/LucasZatta/letter-goxd/internal/lists"
)

func main() {
	startTime := time.Now()
	scraper := lists.New()
	movies := scraper.ScrapeWatchlist("zvttx")
	endTime := time.Now()

	fmt.Printf("EXECUTED IN %v\n", endTime.Sub(startTime))

	for _, m := range *movies {
		fmt.Printf("%v\n", m)
	}
}
