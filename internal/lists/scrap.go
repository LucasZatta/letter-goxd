package lists

import (
	"fmt"
	"os"

	"github.com/anaskhan96/soup"
)

func ScrapeListPreview(username string) *[]MovieDetails {
	movies := []MovieDetails{}
	path := fmt.Sprintf("https://letterboxd.com/%s/watchlist", username)

	for page := 1; ; page++ {
		pathPage := fmt.Sprintf("%s/page/%d", path, page)

		resp, err := soup.Get(pathPage)
		if err != nil {
			os.Exit(1)
		}

		doc := soup.HTMLParse(resp)
		table := doc.Find("ul", "class", "poster-list")
		films := table.FindAll("li", "class", "poster-container")

		for _, filmEntry := range films {
			children := filmEntry.Find("div", "class", "really-lazy-load")

			// movie.ImgSrc = children.Find("img").Attrs()["src"]
			// movie.Name = children.Find("img").Attrs()["alt"]
			moviePath := children.Attrs()["data-target-link"]

			movieDetail := ScrapeMoviePage(moviePath)

			movies = append(movies, movieDetail)
		}

		next := doc.Find("a", "class", "next")
		if next.Error != nil {
			break
		}
	}

	return &movies
}

func ScrapeMoviePage(moviePath string) MovieDetails {
	path := fmt.Sprintf("https://letterboxd.com/%s", moviePath)

	resp, err := soup.Get(path)
	if err != nil {
		os.Exit(1)
	}

	doc := soup.HTMLParse(resp)

}
