package colly

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MK-BK/tk-colly/models"
	"github.com/gocolly/colly/v2"
)

func GetPages() (int, error) {
	var pages int
	c := collector.Clone()

	c.OnHTML(".page_tip", func(e *colly.HTMLElement) {
		reg := regexp.MustCompile(`当前\d+/(\d+)页`)
		match := reg.FindStringSubmatch(e.Text)

		if len(match) == 2 {
			pages, _ = strconv.Atoi(match[1])
		}
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Errorf("colly error: %d, %+v", r.StatusCode, string(r.Body))
	})

	return pages, c.Visit(randomSourceDomain())
}

func GetMovies(page int) ([]*models.Movie, error) {
	movies := make([]*models.Movie, 0)

	c := collector.Clone()
	c.OnHTML(".xing_vb ul li", func(e *colly.HTMLElement) {
		movie := &models.Movie{
			Href:     e.ChildAttr("a", "href"),
			Name:     e.DOM.Find(".xing_vb4").Text(),
			Categoty: e.DOM.Find(".xing_vb5").Text(),
		}
		if ts := e.DOM.Find(".xing_vb7").Text(); ts != "" {
			movie.UpdatedAt, _ = time.Parse(models.TimeFormat, ts)
		}
		if movie.Name != "" && movie.Categoty != "" {
			movies = append(movies, movie)
		}
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Errorf("colly error: %d, %+v", r.StatusCode, string(r.Body))
	})

	if err := c.Visit(strings.Join([]string{randomSourceDomain(), "index/index/page", strconv.Itoa(page)}, "/")); err != nil {
		return nil, err
	}

	if len(movies) > 0 {
		if err := models.GlobalEnvironment.MoviesManager.Save(context.Background(), movies...); err != nil {
			log.Error("Save:", err)
		}
	}

	return movies, nil
}
