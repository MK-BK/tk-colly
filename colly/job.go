package colly

import (
	"time"
)

func Schedule(duration time.Duration) {
	ticker := time.NewTicker(duration)

	for {
		select {
		case <-ticker.C:
			if err := Colly(); err != nil {
				log.Error(err)
			}
		}
	}
}
