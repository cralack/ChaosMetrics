package fetcher

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/utils/rater"
	"go.uber.org/zap"
)

type Fetcher interface {
	rater.RateLimter
	Get(task *Task) ([]byte, error)
}

type BrowserFetcher struct {
	requireRateLimiter bool
	logger             *zap.Logger
	timeout            time.Duration

	rater rater.RateLimter
	pass  chan struct{}
}

var _ Fetcher = &BrowserFetcher{}

func NewBrowserFetcher(opts ...func(*BrowserFetcher)) *BrowserFetcher {
	conf := global.GVA_CONF.Fetcher
	f := defaultFetcher

	for _, opt := range opts {
		opt(f)
	}

	if f.requireRateLimiter {
		if limiter, err := rater.NewSlidingWindowLimiter(
			conf.RateLimiterConfig.Each2Min-2,
			time.Minute*2,
			time.Second/time.Duration(conf.RateLimiterConfig.EachSec),
		); err != nil {
			global.GVA_LOG.Error("rate limiter init failed", zap.Error(err))
			return nil
		} else {
			pass := make(chan struct{})
			f.rater = limiter
			f.pass = pass
			go limiter.StartTimer(pass)
		}
	}
	return f
}

// Get raw json from riot dev port
func (f *BrowserFetcher) Get(task *Task) ([]byte, error) {
	client := &http.Client{
		Timeout: f.timeout * time.Second,
	}
	// set req url
	req, err := http.NewRequestWithContext(
		context.Background(), http.MethodGet, task.URL, nil)
	if err != nil {
		return nil, err
	}

	// set header
	req.Header.Set("User-Agent", task.Header.UserAgent)
	req.Header.Set("Accept-Language", task.Header.AcceptLanguage)
	req.Header.Set("Accept-Charset", task.Header.AcceptCharset)
	req.Header.Set("Origin", task.Header.Origin)
	req.Header.Set("X-Riot-Token", task.Header.ApiToken)

	// require pass signal(rate limiter)
	if f.requireRateLimiter {
		<-f.pass
	}

	// run req
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	defer func() {
		if cerr := resp.Body.Close(); err != nil {
			f.logger.Error("failed to close response body",
				zap.Error(cerr))
		}
	}()

	// get buffer
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		f.logger.Error("read resp body failed",
			zap.Error(err))
		return nil, err
	}

	return body, nil
}

func (f *BrowserFetcher) TryAcquire() bool {
	return f.rater.TryAcquire()
}
