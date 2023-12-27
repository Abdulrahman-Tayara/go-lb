package buckettokens

import (
	"time"
)

func startResetBackground(cfg *Configs, data *buckets) {
	go func() {
		for {
			reset(cfg.Rate, data)
			time.Sleep(cfg.Rate)
		}
	}()
}

func reset(rate time.Duration, data *buckets) {
	data.deleteFunc(func(b *bucket) bool {
		return b.shouldReset(rate)
	})
}
