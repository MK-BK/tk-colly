package colly

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/MK-BK/tk-colly/models"

	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
)

const (
	domain     = "https://tkzy2.com"
	pageURL    = "index/index/page"
	timeFormat = "2006-01-02 15:04:05"
)

var (
	log          = logrus.New()
	collector    *colly.Collector
	concurrency  = 4
	latestUpdate time.Time
)

func init() {
	collector = colly.NewCollector()

	// collector.SetClient(&http.Client{
	// 	Transport: &http.Transport{
	// 		Proxy: func(r *http.Request) (*url.URL, error) {
	// 			proxyURL, err := url.Parse("http://157.100.7.146:999")
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			return proxyURL, nil
	// 		},
	// 	},
	// })
}

func Colly() error {
	count, err := GetCount()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	channels := make(chan int, concurrency)
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				i, ok := <-channels
				if !ok {
					return
				}

				CollectorPage(i)
			}
		}()
	}

	for i := 1; i <= count; i++ {
		channels <- i
	}

	close(channels)
	wg.Wait()
	return nil
}

func GetCount() (int, error) {
	var count int

	c := collector.Clone()
	c.OnHTML(".page_tip", func(e *colly.HTMLElement) {
		reg := regexp.MustCompile(`当前\d+/(\d+)页`)
		match := reg.FindStringSubmatch(e.Text)
		if len(match) == 2 {
			count, _ = strconv.Atoi(match[1])
		}
	})

	c.Visit(domain)

	return count, nil
}

func CollectorPage(pageIndex int) {
	movies := make([]*models.Movie, 0)

	now := time.Now()
	defer func() {
		log.Infof("Visit %s, spend time: %+v", fmt.Sprintf("%s/%s/%d/", domain, pageURL, pageIndex), time.Now().Sub(now))
		if len(movies) > 0 {
			if err := models.GlobalEnvironment.MovieInterface.Save(context.Background(), movies...); err != nil {
				log.Error("Save err:", err)
			}
		}
	}()

	c := collector.Clone()
	c.OnHTML(".xing_vb ul li", func(e *colly.HTMLElement) {
		movie := &models.Movie{
			Href: e.ChildAttr("a", "href"),
		}

		movie.Name = e.DOM.Find(".xing_vb4").Text()
		movie.Categoty = e.DOM.Find(".xing_vb5").Text()

		if e.DOM.Find(".xing_vb7").Text() != "" {
			if update_at, err := time.Parse(timeFormat, e.DOM.Find(".xing_vb7").Text()); err != nil {
				log.Errorf("collectorPage %s, err: %+v", fmt.Sprintf("%s/%s/%d/", domain, pageURL, pageIndex), err)
			} else {
				movie.UpdatedAt = update_at
			}
		}

		if movie.Name != "" && movie.Categoty != "" {
			movies = append(movies, movie)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Errorf("Request URL: %s, Failed with error: %+v", r.Request.URL, r.StatusCode)
	})

	c.Visit(fmt.Sprintf("%s/%s/%d/", domain, pageURL, pageIndex))

	collectorPageView(movies...)
}

func collectorPageView(movies ...*models.Movie) error {
	movieViews := make([]*models.MovieView, 0)

	defer func() {
		if len(movieViews) > 0 {
			if err := models.GlobalEnvironment.MovieInterface.SaveView(context.Background(), movieViews...); err != nil {
				log.Error("SaveView err:", err)
			}
		}
	}()

	var wg sync.WaitGroup

	for i, movie := range movies {
		if !latestUpdate.IsZero() && latestUpdate.Sub(movie.UpdatedAt) > 0 {
			continue
		}

		wg.Add(1)

		go func(index int) {
			defer func() {
				wg.Done()
			}()

			c := collector.Clone()

			movieView := &models.MovieView{}

			c.OnHTML(".vodinfobox ul", func(e *colly.HTMLElement) {
				movieView.DisplayName = e.DOM.Find("li:nth-child(1) span").Text()
				movieView.Director = e.DOM.Find("li:nth-child(2) span").Text()
				movieView.Actors = e.DOM.Find("li:nth-child(3) span").Text()
				movieView.Category = e.DOM.Find("li:nth-child(4) span").Text()
				movieView.Region = e.DOM.Find("li:nth-child(5) span").Text()
				movieView.Language = e.DOM.Find("li:nth-child(6) span").Text()
				movieView.Released = e.DOM.Find("li:nth-child(7) span").Text()

				if e.DOM.Find("li:nth-child(8) span").Text() != "" {
					if updateAt, err := time.Parse(timeFormat, e.DOM.Find("li:nth-child(8) span").Text()); err != nil {
						log.Errorf("collectorPageView %s, err: %+v", fmt.Sprintf("%s%s", domain, movies[index].Href), err)
					} else {
						movieView.UpdatedAt = updateAt
					}
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

			c.OnError(func(r *colly.Response, err error) {
				log.Errorf("Request URL: %s, Failed with error: %+v", r.Request.URL, r.StatusCode)
			})

			c.Visit(fmt.Sprintf("%s%s", domain, movies[index].Href))

			movieViews = append(movieViews, movieView)
		}(i)
	}

	wg.Wait()
	return nil
}

func GetPageMovie(pageIndex int) ([]*models.Movie, error) {
	movies := make([]*models.Movie, 0)

	defer func() {
		if len(movies) > 0 {
			if err := models.GlobalEnvironment.MovieInterface.Save(context.Background(), movies...); err != nil {
				fmt.Println("save error:", err)
			}
		}
	}()

	c := collector.Clone()

	c.OnHTML(".xing_vb ul li", func(e *colly.HTMLElement) {
		movie := &models.Movie{
			Href: e.ChildAttr("a", "href"),
		}

		movie.Name = e.DOM.Find(".xing_vb4").Text()
		movie.Categoty = e.DOM.Find(".xing_vb5").Text()

		if e.DOM.Find(".xing_vb7").Text() != "" {
			if update_at, err := time.Parse(timeFormat, e.DOM.Find(".xing_vb7").Text()); err != nil {
				log.Errorf("collectorPage %s, err: %+v", fmt.Sprintf("%s/%s/%d/", domain, pageURL, pageIndex), err)
			} else {
				movie.UpdatedAt = update_at
			}
		}

		if movie.Name != "" && movie.Categoty != "" {
			movies = append(movies, movie)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Errorf("Request URL: %s, Failed with error: %+v", r.Request.URL, r.StatusCode)
	})

	c.Visit(fmt.Sprintf("%s/%s/%d/", domain, pageURL, pageIndex))

	return movies, nil
}

func SetLatestUpdate(latest time.Time) {
	latestUpdate = latest
}

func GetLatestUpdate() time.Time {
	return latestUpdate
}
