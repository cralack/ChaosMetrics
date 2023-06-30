package test

import (
	_ "github.com/cralack/ChaosMetrics/server/init"
	"github.com/cralack/ChaosMetrics/server/pkg/fetcher"
)

var f fetcher.Fetcher
var path string

func init() {
	f = fetcher.NewBrowserFetcher()
	path = "./local_json/"
}
