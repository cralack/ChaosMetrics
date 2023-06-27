package fetcher

import (
	"context"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/cralack/ChaosMetrics/server/config"
	"github.com/cralack/ChaosMetrics/server/global"

	"go.uber.org/zap"
)

type Fetcher interface {
	Get(url string) ([]byte, error)
}

type BrowserFetcher struct {
	Conf    *config.FetcherConfig
	Logger  *zap.Logger
	Timeout time.Duration
}

var _ Fetcher = &BrowserFetcher{}

func NewBrowserFetcher() *BrowserFetcher {
	conf := global.GVA_CONF.Fetcher
	if conf.HeaderConfig.XRiotToken == "" {
		workDir := global.GVA_CONF.DirTree.WorkDir
		path := filepath.Join(workDir, "pkg", "fetcher")
		buff, err := ioutil.ReadFile(path + "/api_key")
		if err != nil {
			global.GVA_LOG.Error("get api key failed",
				zap.Error(err))
		}
		conf.HeaderConfig.XRiotToken = string(buff)
	}
	return &BrowserFetcher{
		Timeout: global.GVA_CONF.Fetcher.Timeout,
		Logger:  global.GVA_LOG,
		Conf:    conf,
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
