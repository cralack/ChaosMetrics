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
	logger:             global.GVA_LOG,
	timeout:            global.GVA_CONF.Fetcher.Timeout,
	header: &Header{
		AcceptLanguage: global.GVA_CONF.Fetcher.HeaderConfig.AcceptLanguage,
		AcceptCharset:  global.GVA_CONF.Fetcher.HeaderConfig.AcceptCharset,
		Origin:         global.GVA_CONF.Fetcher.HeaderConfig.Origin,
		UserAgent:      global.GVA_CONF.Fetcher.HeaderConfig.UserAgent,
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
