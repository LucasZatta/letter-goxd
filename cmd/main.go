package main

import (
	"fmt"
	"time"

	"github.com/LucasZatta/letter-goxd/internal/lists"
)

func main() {
	startTime := time.Now()
	scraper := lists.New()
	movies := scraper.PerformanceTest("zvttx")
	endTime := time.Now()

	fmt.Printf("EXECUTED IN %v\n", endTime.Sub(startTime))

	for i, m := range *movies {
		fmt.Printf("%v\n", m)
		fmt.Println(i)
	}
}
