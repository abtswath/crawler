package crawler

import "net/url"

type Options struct {
	UserAgent string
	PoolSize  int
	Target    *url.URL
}
