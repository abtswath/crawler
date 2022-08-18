package browser

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
	Cookies      []*proto.NetworkCookieParam
	MaxPageCount int
	UserAgent    string
	Headers      []string
}

type PageOption struct {
	Timeout time.Duration
}
