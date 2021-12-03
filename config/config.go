package config

import "net/url"

type Options struct {
	BrowserPath    string
	Incognito      bool
	Headless       bool
	Proxy          string
	Headers        map[string]string
	PoolSize       int
	Target         *url.URL
	IgnoreKeywords []string
}
