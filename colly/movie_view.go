package colly

import (
	"context"
	"fmt"
	"time"

	"github.com/MK-BK/tk-colly/models"
	"github.com/gocolly/colly/v2"
)

func GetMoviesView(ctx context.Context, movie *models.Movie) error {
	if !latestUpdate.IsZero() && latestUpdate.Sub(movie.UpdatedAt) > 0 {
		return nil
	}

	log.Info("collector movie:", movie.Categoty, movie.Href, movie.Name)
	c := collector.Clone()

	movieView := &models.MovieView{
		MovieID: int(movie.ID),
	}

	c.OnHTML(".vodImg", func(e *colly.HTMLElement) {
		movieView.ImagePath = e.ChildAttr("img", "src")
	})

	c.OnHTML(".vodInfo .vodh", func(e *colly.HTMLElement) {
		movieView.Name = e.ChildText("h2")
		movieView.Rating = e.ChildText("label")
	})

	c.OnHTML(".warp", func(e *colly.HTMLElement) {
		e.ForEach(".ibox.playBox", func(i int, el *colly.HTMLElement) {
			if i == 0 {
				movieView.Description = el.ChildText(".vodplayinfo")
			}
		})
	})

	c.OnHTML(".vodinfobox ul", func(e *colly.HTMLElement) {
		movieView.DisplayName = e.DOM.Find("li:nth-child(1) span").Text()
		movieView.Director = e.DOM.Find("li:nth-child(2) span").Text()
		movieView.Actors = e.DOM.Find("li:nth-child(3) span").Text()
		movieView.Category = e.DOM.Find("li:nth-child(4) span").Text()
		movieView.Region = e.DOM.Find("li:nth-child(5) span").Text()
		movieView.Language = e.DOM.Find("li:nth-child(6) span").Text()
		movieView.Released = e.DOM.Find("li:nth-child(7) span").Text()
		if ts := e.DOM.Find("li:nth-child(8) span").Text(); ts != "" {
			movieView.UpdatedAt, _ = time.Parse(models.TimeFormat, ts)
		}
	})

	c.OnHTML("div.ibox.playBox", func(e *colly.HTMLElement) {
		e.ForEach("div.vodplayinfo", func(_ int, el *colly.HTMLElement) {
			play := &models.Play{
				Resource: el.ChildText("h3"),
			}

			el.ForEach("ul li", func(_ int, li *colly.HTMLElement) {
				play.Name = li.ChildAttr("a", "href")
				play.URL = li.ChildText("a")
			})
			if movieView.PlayLists == nil {
				movieView.PlayLists = make([]*models.Play, 0)
			}
			movieView.PlayLists = append(movieView.PlayLists, play)
		})
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Errorf("colly error: %d, %+v", r.StatusCode, string(r.Body))
	})

	if err := c.Visit(fmt.Sprintf("%s%s", randomSourceDomain(), movie.Href)); err != nil {
		return err
	}

	return models.GlobalEnvironment.MoviesManager.SaveView(ctx, movieView)
}
