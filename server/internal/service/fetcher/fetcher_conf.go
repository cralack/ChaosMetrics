package fetcher

import (
	"github.com/cralack/ChaosMetrics/server/utils"
)

type Header struct {
	AcceptLanguage string
	AcceptCharset  string
	ApiToken       string
	Origin         string
	UserAgent      string
}

func WithRateLimiter(flag bool) func(*BrowserFetcher) {
	return func(fetcher *BrowserFetcher) {
		if !flag {
			fetcher.requireRateLimiter = false
			fetcher.rater = nil
		}
	}
}

func WithHeader(h *Header) func(*BrowserFetcher) {
	return func(fetcher *BrowserFetcher) {
		if h != nil {
			fetcher.header = h
		}
	}
}

func WithRandomUA() func(*BrowserFetcher) {
	return func(fetcher *BrowserFetcher) {
		fetcher.header.UserAgent = utils.GenerateRandomUA()
	}
}

func WithToken(token string) func(*BrowserFetcher) {
	return func(fetcher *BrowserFetcher) {
		fetcher.header.ApiToken = token
	}
}
