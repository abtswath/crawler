package crawler

import (
	"log"
	"time"

	"github.com/go-rod/rod/lib/proto"
)

type Option struct {
	Logger       *log.Logger
	Headless     bool
	Bin          string
	Proxy        string
	Timeout      time.Duration
	Cookies      []*proto.NetworkCookieParam
	MaxPageCount int
	UserAgent    string
	Headers      []string
	PageTimeout  time.Duration
	Exclusions   []string
}

type SetOptionValueFunc = func(*Option)

func WithLogger(logger *log.Logger) SetOptionValueFunc {
	return func(o *Option) {
		if o.Logger == nil {
			o.Logger = logger
		}
	}
}

func WithHeadless(headless bool) SetOptionValueFunc {
	return func(o *Option) {
		if !o.Headless {
			o.Headless = headless
		}
	}
}

func WithBin(bin string) SetOptionValueFunc {
	return func(o *Option) {
		if o.Bin == "" {
			o.Bin = bin
		}
	}
}

func WithProxy(proxy string) SetOptionValueFunc {
	return func(o *Option) {
		if o.Proxy == "" {
			o.Proxy = proxy
		}
	}
}

func WithTimeout(timeout time.Duration) SetOptionValueFunc {
	return func(o *Option) {
		if o.Timeout == 0 {
			o.Timeout = timeout
		}
	}
}

func WithCookies(cookies []*proto.NetworkCookieParam) SetOptionValueFunc {
	return func(o *Option) {
		if o.Cookies == nil {
			o.Cookies = cookies
		}
	}
}

func WithMaxPageCount(maxPageCount int) SetOptionValueFunc {
	return func(o *Option) {
		if o.MaxPageCount == 0 {
			o.MaxPageCount = maxPageCount
		}
	}
}

func WithUserAgent(userAgent string) SetOptionValueFunc {
	return func(o *Option) {
		if o.UserAgent == "" {
			o.UserAgent = userAgent
		}
	}
}

func WithHeaders(headers []string) SetOptionValueFunc {
	return func(o *Option) {
		if o.Headers == nil {
			o.Headers = headers
		}
	}
}

func WithPageTimeout(timeout time.Duration) SetOptionValueFunc {
	return func(o *Option) {
		if o.PageTimeout == 0 {
			o.PageTimeout = timeout
		}
	}
}
