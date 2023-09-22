package colly

import (
	"time"

	"github.com/robfig/cron"
)

func Schedule(duration time.Duration) {
	c := cron.New()

	c.AddFunc("0 0 12 * * ?", func() {
		now := time.Now()
		collector := NewCollector(&now)

		if err := collector.Colly(); err != nil {
			collector.Logger.Error(err)
		}
	})
	c.Start()
}
