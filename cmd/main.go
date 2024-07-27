package main

import (
	"fmt"
	"os"

	"github.com/anaskhan96/soup"
)

type MoviePreview struct {
	Name string
	Path string
}

func main() {
	movies := []MoviePreview{}

	for page := 1; ; page++ {
		resp, err := soup.Get(fmt.Sprintf("https://letterboxd.com/zvttx/watchlist/page/%d", page))
		if err != nil {
			os.Exit(1)
		}
		doc := soup.HTMLParse(resp)
		table := doc.Find("ul", "class", "poster-list")
		filmes := table.FindAll("li", "class", "poster-container")

		for _, filmEntry := range filmes {
			movie := MoviePreview{}
			children := filmEntry.Find("div", "class", "really-lazy-load")

			movie.Name = children.Find("img").Attrs()["alt"]
			movie.Path = children.Attrs()["data-target-link"]

			movies = append(movies, movie)
		}
		fmt.Println(page)
		next := doc.Find("a", "class", "next")
		if next.Error != nil {
			break
		}
	}

	fmt.Printf("%+v\n", movies)

}
