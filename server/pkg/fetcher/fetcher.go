package fetcher

import (
	"context"
	"github.com/cralack/ChaosMetrics/server/global"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Fetcher interface {
	Get(url string) ([]byte, error)
}

type BrowserFetcher struct {
	Timeout time.Duration
	Logger  *zap.Logger
}

var _ Fetcher = &BrowserFetcher{}

func NewBrowserFetcher() *BrowserFetcher {
	return &BrowserFetcher{
		Timeout: global.GVA_CONF.Fetcher.Timeout,
		Logger:  global.GVA_LOG,
	}
}

// get raw json from riot dev port
func (f *BrowserFetcher) Get(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: f.Timeout * time.Second,
	}
	//set req url
	req, err := http.NewRequestWithContext(
		context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	//set header
	header := global.GVA_CONF.Fetcher.HeaderConfig
	req.Header.Set("User-Agent", header.UserAgent)
	req.Header.Set("Accept-Language", header.AcceptLanguage)
	req.Header.Set("Accept-Charset", header.AcceptCharset)
	req.Header.Set("Origin", header.Origin)
	req.Header.Set("X-Riot-Token", header.XRiotToken)

	//run req
	resp, err := client.Do(req)
	if err != nil {
		f.Logger.Error("fetch failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	//get buffer
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		f.Logger.Error("read resp body failed", zap.Error(err))
		return nil, err
	}

	return body, nil
}
