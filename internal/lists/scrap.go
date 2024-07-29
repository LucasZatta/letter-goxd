package lists

import (
	"fmt"
	"os"
	"strconv"

	"github.com/LucasZatta/letter-goxd/internal/util"
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
	var movieDetails MovieDetails
	path := fmt.Sprintf("https://letterboxd.com/%s", moviePath)

	resp, err := soup.Get(path)
	if err != nil {
		os.Exit(1)
	}

	doc := soup.HTMLParse(resp)

	movieDetails.Url = doc.Find("meta", "property", "og:url").Attrs()["content"]
	movieDetails.Name = doc.Find("meta", "property", "og:title").Attrs()["content"]
	movieDetails.Description = doc.Find("meta", "property", "og:description").Attrs()["content"]
	movieDetails.Image = doc.Find("meta", "property", "og:image").Attrs()["content"]

	duration, _ := strconv.Atoi(util.ClearString(doc.Find("p", "class", "text-link").Text()))
	movieDetails.Duration = duration

	movieDetails.Director = doc.Find("meta", "name", "twitter:data1").Attrs()["content"]
	rating := doc.Find("meta", "name", "twitter:data2")
	if rating.Error == nil {
		movieDetails.Rating = doc.Find("meta", "name", "twitter:data2").Attrs()["content"]
	} //treat fields for missing entry error

	return MovieDetails{}

}
