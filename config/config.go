package config

import (
	"net/url"
	"time"
)

type Options struct {
	Timeout        time.Duration
	BrowserPath    string
	Incognito      bool
	Headless       bool
	Proxy          string
	Headers        map[string]string
	PoolSize       int
	Target         *url.URL
	PageTimeout    time.Duration
	IgnoreKeywords []string
	UploadFile     string
}
