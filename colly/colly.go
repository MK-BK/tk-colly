package colly

import (
	"context"
	"net/http"
	"time"

	"github.com/MK-BK/tk-colly/common"
	"github.com/MK-BK/tk-colly/models"
	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
)

type Collector struct {
	Concurrency  int
	Collector    *colly.Collector
	Logger       *logrus.Logger
	ScheduleTime *time.Time
	UpdateTime   time.Time
}

func NewCollector(scheduleTime *time.Time) *Collector {
	c := colly.NewCollector()
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
	})

	return &Collector{
		Concurrency:  20,
		Logger:       logrus.New(),
		Collector:    c,
		ScheduleTime: scheduleTime,
	}
}

func (c *Collector) Colly() error {
	pages, err := c.GetPages()
	if err != nil {
		return err
	}

	chs := make(chan *models.Movie, c.Concurrency)
	for i := 0; i < c.Concurrency; i++ {
		go c.Process(chs)
	}

	for i := 1; i <= pages; i++ {
		movies, err := c.GetMovies(i)
		if err != nil {
			c.Logger.Error(err)
			continue
		}

		for _, movie := range movies {
			if c.ScheduleTime != nil && movie.UpdatedAt.Sub(*c.ScheduleTime) < 0 {
				goto finished
			}

			chs <- movie
		}
	}

finished:
	close(chs)
	return nil
}

func (c *Collector) Process(chs chan *models.Movie) {
	for {
		movie, ok := <-chs
		if !ok {
			return
		}

		movieView, moviePlayer, err := c.GetMoviesView(context.Background(), movie)
		if err != nil {
			c.Logger.Error(err)
			continue
		}

		handler := func(movieView *models.MovieView, player *models.MoviePlayer) error {
			tx := common.DB.Begin()

			if err := models.GlobalEnvironment.MoviesManager.CreateMovieView(tx, movieView); err != nil {
				return tx.Rollback().Error
			}

			moviePlayer.MovieID = int(movieView.ID)

			if err := models.GlobalEnvironment.MoviesManager.CreateMoviePlayer(tx, moviePlayer); err != nil {
				return tx.Rollback().Error
			}

			return tx.Commit().Error
		}

		if err := handler(movieView, moviePlayer); err != nil {
			c.Logger.WithField("href", movieView.Href).Info()
			c.Logger.Error(err)
		}
	}
}

func (c *Collector) SetLatestUpdate(updateTime time.Time) {
	c.UpdateTime = updateTime
}

func (c *Collector) GetLatestUpdate() time.Time {
	return c.UpdateTime
}
