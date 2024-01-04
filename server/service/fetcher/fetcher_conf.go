package fetcher

import (
	"github.com/cralack/ChaosMetrics/server/internal/global"
)

var defaultFetcher = &BrowserFetcher{
	requireRateLimiter: false,
	logger:             global.GVA_LOG,
	timeout:            global.GVA_CONF.Fetcher.Timeout,
}

func WithRateLimiter(flag bool) func(*BrowserFetcher) {
	return func(fetcher *BrowserFetcher) {
		if !flag {
			fetcher.requireRateLimiter = false
			fetcher.rater = nil
		}
	}
}
