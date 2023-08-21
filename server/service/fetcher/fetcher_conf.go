package fetcher

import (
	"github.com/cralack/ChaosMetrics/server/config"
)

type Option func(conf *config.FetcherConfig)

func WithAPIToken(token string) Option {
	return func(conf *config.FetcherConfig) {
		conf.HeaderConfig.XRiotToken = token
	}
}

func WithRateLimiter(flag bool) Option {
	return func(conf *config.FetcherConfig) {
		conf.RequireRateLimiter = flag
	}
}