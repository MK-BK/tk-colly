package job

import (
	"time"

	"github.com/MK-BK/tk-colly/colly"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Schedule(duration time.Duration) {
	ticker := time.NewTicker(duration)

	for {
		select {
		case <-ticker.C:
			pageCount, err := colly.GetCount()
			if err != nil {
				log.Error(err)
				return
			}

			for i := 1; i <= pageCount; i++ {
				movices, err := colly.GetPageMovie(i)
				if err != nil {
					log.Error("Job Schedule:", err)
					return
				}

				if len(movices) > 0 && !colly.GetLatestUpdate().IsZero() && colly.GetLatestUpdate().Sub(movices[0].UpdatedAt) > 0 {
					break
				}

				colly.CollectorPage(i)
			}
		}
	}
}
