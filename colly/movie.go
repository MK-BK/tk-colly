package colly

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MK-BK/tk-colly/models"

	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
)

func (c *Collector) GetPages() (int, error) {
	var pages int

	collector := c.Collector.Clone()

	collector.OnHTML(".page_tip", func(e *colly.HTMLElement) {
		reg := regexp.MustCompile(`当前\d+/(\d+)页`)
		match := reg.FindStringSubmatch(e.Text)

		if len(match) == 2 {
			pages, _ = strconv.Atoi(match[1])
		}
	})

	if err := collector.Visit(randomSourceDomain()); err != nil {
		return pages, err
	}

	return pages, nil
}

func (c *Collector) GetMovies(page int) ([]*models.Movie, error) {
	movies := make([]*models.Movie, 0)

	collector := c.Collector.Clone()

	collector.OnHTML(".xing_vb ul li", func(e *colly.HTMLElement) {
		movie := &models.Movie{
			Href:     e.ChildAttr("a", "href"),
			Name:     e.DOM.Find(".xing_vb4").Text(),
			Categoty: e.DOM.Find(".xing_vb5").Text(),
		}

		if ts := e.DOM.Find(".xing_vb7").Text(); ts != "" {
			movie.UpdatedAt, _ = time.Parse(TimeFormat, ts)
		}

		if movie.Name != "" && movie.Categoty != "" {
			movies = append(movies, movie)
		}
	})

	url := strings.Join([]string{randomSourceDomain(), "index/index/page", strconv.Itoa(page)}, "/")
	if err := collector.Visit(url); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("get movie page error: %s", url))
	}

	return movies, nil
}

func (c *Collector) GetMoviesView(ctx context.Context, movie *models.Movie) (*models.MovieView, *models.MoviePlayer, error) {
	collector := c.Collector.Clone()

	movieView := &models.MovieView{
		Href:     movie.Href,
		Category: movie.Categoty,
	}

	collector.OnHTML(".vodImg", func(e *colly.HTMLElement) {
		movieView.ImagePath = e.ChildAttr("img", "src")
	})

	collector.OnHTML(".vodInfo .vodh", func(e *colly.HTMLElement) {
		movieView.Name = e.ChildText("h2")
		movieView.Rating = e.ChildText("label")
	})

	collector.OnHTML(".warp", func(e *colly.HTMLElement) {
		e.ForEach(".ibox.playBox", func(i int, el *colly.HTMLElement) {
			if i == 0 {
				movieView.Description = el.ChildText(".vodplayinfo")
			}
		})
	})

	collector.OnHTML(".vodinfobox ul", func(e *colly.HTMLElement) {
		movieView.DisplayName = e.DOM.Find("li:nth-child(1) span").Text()
		movieView.Director = e.DOM.Find("li:nth-child(2) span").Text()
		movieView.Actors = e.DOM.Find("li:nth-child(3) span").Text()
		movieView.SubCategory = e.DOM.Find("li:nth-child(4) span").Text()
		movieView.Region = e.DOM.Find("li:nth-child(5) span").Text()
		movieView.Language = e.DOM.Find("li:nth-child(6) span").Text()
		movieView.Released = e.DOM.Find("li:nth-child(7) span").Text()

		if ts := e.DOM.Find("li:nth-child(8) span").Text(); ts != "" {
			movieView.UpdatedAt, _ = time.Parse(TimeFormat, ts)
		}
	})

	moviePlayer := &models.MoviePlayer{
		Players: make([]*models.Player, 0),
	}

	collector.OnHTML("font > .vodplayinfo", func(e *colly.HTMLElement) {
		if strings.Contains(e.ChildText("h3"), "tkyun") {
			e.ForEach("ul li", func(_ int, li *colly.HTMLElement) {
				moviePlayer.Players = append(moviePlayer.Players, &models.Player{
					Name: li.ChildAttr("a", "href"),
					URL:  li.ChildText("a"),
				})
			})
		}
	})

	url := fmt.Sprintf("%s%s", randomSourceDomain(), movie.Href)
	if err := collector.Visit(url); err != nil {
		return nil, nil, errors.WithMessage(err, fmt.Sprintf("get movie view error: %s", url))
	}

	return movieView, moviePlayer, nil
}
