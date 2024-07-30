package lists

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/LucasZatta/letter-goxd/internal/util"
	"github.com/anaskhan96/soup"
	"github.com/cenkalti/backoff/v4"
)

type Scrape interface {
	ScrapeWatchlist(username string) *[]MovieDetails
	ScrapeMoviePage(moviePath string) (*MovieDetails, error)
}

type scrape struct {
}

func New() *scrape {
	return &scrape{}
}

func (s *scrape) ScrapeWatchlist(username string) *[]MovieDetails {
	links := make([]string, 0)
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

			links = append(links, moviePath)
		}
		next := doc.Find("a", "class", "next")
		if next.Error != nil {
			break
		}
	}

	movies := make([]MovieDetails, len(links))

	wg := sync.WaitGroup{}
	for i, link := range links {
		wg.Add(1)
		go func(link string) {
			movie, err := s.ScrapeMoviePage(link)
			if err != nil {
				log.Fatal(err)
			}
			movies[i] = *movie
			wg.Done()
		}(link)

	}
	wg.Wait()

	return &movies
}

func (s *scrape) ScrapeMoviePage(moviePath string) (*MovieDetails, error) {
	var movieDetails MovieDetails

	operation := func() error {
		path := fmt.Sprintf("https://letterboxd.com%s", moviePath)
		movieDetails.Url = path

		resp, err := soup.Get(path)
		if err != nil {
			os.Exit(1)
		}

		doc := soup.HTMLParse(resp)

		//implement retry with backoff
		nameRoot := doc.Find("meta", "property", "og:title")
		if nameRoot.Error != nil {
			//should return err and bail
			log.Print("cannot find movie name")
			return nameRoot.Error
		}
		movieDetails.Name = nameRoot.Attrs()["content"]

		descriptionRoot := doc.Find("meta", "property", "og:description")
		if descriptionRoot.Error != nil {
			log.Print("cannot find movie description")
			return descriptionRoot.Error
		} else {
			movieDetails.Description = descriptionRoot.Attrs()["content"]
		}

		imageRoot := doc.Find("meta", "property", "og:image")
		if imageRoot.Error != nil {
			log.Print("cannot find movie image")
			return imageRoot.Error
		} else {
			movieDetails.Image = imageRoot.Attrs()["content"]
		}

		durationRoot := doc.Find("p", "class", "text-link")
		if durationRoot.Error != nil {
			log.Print("cannot find movie duration")
			return durationRoot.Error
		} else {
			duration := util.StringElementToInt(durationRoot.Text())
			movieDetails.Duration = duration
		}

		directorsRoot := doc.Find("meta", "name", "twitter:data1")
		if directorsRoot.Error != nil {
			log.Print("cannot find directors")
			return directorsRoot.Error
		} else {
			movieDetails.Director = directorsRoot.Attrs()["content"]
		}

		ratingRoot := doc.Find("meta", "name", "twitter:data2")
		if ratingRoot.Error != nil {
			log.Print("couldnt find movie rating")
			// return ratingRoot.Error
		} else {
			movieDetails.Rating = ratingRoot.Attrs()["content"]
		}
		return nil
	}

	err := backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &movieDetails, nil
}
