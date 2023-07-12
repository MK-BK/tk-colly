package colly

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/MK-BK/tk-colly/models"

	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
)

var once sync.Once
var collector *colly.Collector
var latestUpdate time.Time

var log = logrus.New()

var concurrency = 6

var sourceDomains = []string{
	"https://tkzy1.com",
	// "https://tkzy2.com",
	// "https://tkzy3.com",
	// "https://tkzy4.com",
	// "https://tkzy5.com",
	// "https://tkzy6.com",
	// "https://tkzy7.com",
	// "https://tkzy8.com",
	// "https://tkzy9.com",
}

func randomSourceDomain() string {
	randInt := rand.Intn(len(sourceDomains))
	return sourceDomains[randInt]
}

func init() {
	once.Do(func() {
		collector = colly.NewCollector(
			colly.TraceHTTP(),
		)

		collector.SetProxy("http://127.0.0.1:58591")
	})
}

func Colly() error {
	pages, err := GetPages()
	if err != nil {
		return err
	}

	movieChans := make(chan *models.Movie, concurrency)
	for i := 0; i < concurrency; i++ {
		go worker(movieChans)
	}

	for i := 1; i < pages; i++ {
		movies, err := GetMovies(i)
		if err != nil {
			log.Error("getMOvies error:", err)
			continue
		}
		for _, movie := range movies {
			movieChans <- movie
		}
	}

	close(movieChans)
	return nil
}

func worker(movieChans chan *models.Movie) {
	for {
		movie, ok := <-movieChans
		if !ok {
			return
		}
		if err := GetMoviesView(context.Background(), movie); err != nil {
			log.Error(err)
		}
	}
}

func SetLatestUpdate(latest time.Time) {
	latestUpdate = latest
}

func GetLatestUpdate() time.Time {
	return latestUpdate
}
