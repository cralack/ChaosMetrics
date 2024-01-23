package fetcher

import (
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/utils"
)

type Header struct {
	AcceptLanguage string
	AcceptCharset  string
	ApiToken       string
	Origin         string
	UserAgent      string
}

var defaultFetcher = &BrowserFetcher{
	requireRateLimiter: false,
	logger:             global.GvaLog,
	timeout:            global.GvaConf.Fetcher.Timeout,
	header: &Header{
		AcceptLanguage: global.GvaConf.Fetcher.HeaderConfig.AcceptLanguage,
		AcceptCharset:  global.GvaConf.Fetcher.HeaderConfig.AcceptCharset,
		Origin:         global.GvaConf.Fetcher.HeaderConfig.Origin,
		UserAgent:      global.GvaConf.Fetcher.HeaderConfig.UserAgent,
	},
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
