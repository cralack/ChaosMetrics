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
	requireRateLimiter: true,
	logger:             global.ChaLogger,
	timeout:            global.ChaConf.Fetcher.Timeout,
	header: &Header{
		AcceptLanguage: global.ChaConf.Fetcher.HeaderConfig.AcceptLanguage,
		AcceptCharset:  global.ChaConf.Fetcher.HeaderConfig.AcceptCharset,
		Origin:         global.ChaConf.Fetcher.HeaderConfig.Origin,
		UserAgent:      global.ChaConf.Fetcher.HeaderConfig.UserAgent,
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
