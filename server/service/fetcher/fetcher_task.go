package fetcher

import (
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/utils"
)

type Task struct {
	URL    string
	Header *Header
}

func NewTask(opts ...func(*Task)) *Task {
	t := &Task{
		Header: defaultHeader,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

type Header struct {
	AcceptLanguage string
	AcceptCharset  string
	ApiToken       string
	Origin         string
	UserAgent      string
}

var defaultHeader = &Header{
	AcceptLanguage: global.GVA_CONF.Fetcher.HeaderConfig.AcceptLanguage,
	AcceptCharset:  global.GVA_CONF.Fetcher.HeaderConfig.AcceptCharset,
	Origin:         global.GVA_CONF.Fetcher.HeaderConfig.Origin,
	UserAgent:      global.GVA_CONF.Fetcher.HeaderConfig.UserAgent,
}

func WithURL(url string) func(*Task) {
	return func(t *Task) {
		t.URL = url
	}
}

func WithToken(apitoken string) func(*Task) {
	return func(t *Task) {
		t.Header.ApiToken = apitoken
	}
}

func WithRandomUA() func(*Task) {
	return func(t *Task) {
		t.Header.UserAgent = utils.GenerateRandomUA()
	}
}
