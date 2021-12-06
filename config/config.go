package config

import (
	"net/url"
	"time"
)

type Option struct {
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

func NewOption(target *url.URL) *Option {
	return &Option{
		Timeout:     time.Minute * 5,
		Incognito:   true,
		Headless:    true,
		Headers:     map[string]string{},
		PoolSize:    20,
		Target:      target,
		PageTimeout: time.Second * 5,
	}
}
