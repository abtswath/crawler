package crawler

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/defaults"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/sirupsen/logrus"
	"time"
)

type browserOption struct {
	bin         string
	incognito   bool
	headless    bool
	proxy       string
	pageTimeout time.Duration
	logger      *logrus.Logger
	trace       bool
	cookies     []*proto.NetworkCookieParam
}

func newBrowser(opts browserOption) (*rod.Browser, error) {
	defaults.Show = true
	defaults.Devtools = true
	defaults.Slow = time.Second
	l := launcher.New().
		Leakless(true).
		Headless(opts.headless).
		NoSandbox(true).
		Set("disable-gpu").
		Set("disable-web-security").
		Set("disable-xss-auditor").
		Set("no-sandbox").
		Set("disable-setuid-sandbox").
		Set("allow-running-insecure-content").
		Set("disable-popup-blocking").
		Set("disable-webgl").
		Set("ignore-certificate-errors").
		Set("disable-popup-blocking").
		Set("disable-images")
	if opts.incognito {
		l.Set("incognito")
	}
	if opts.proxy != "" {
		l.Proxy(opts.proxy)
	}
	if opts.bin != "" {
		l.Bin(opts.bin)
	}
	controlURL, err := l.Launch()
	if err != nil {
		return nil, err
	}

	browser := rod.New().
		ControlURL(controlURL)
	err = browser.Connect()
	if err != nil {
		_ = browser.Close()
		return nil, err
	}
	err = browser.SetCookies(opts.cookies)
	if err != nil {
		_ = browser.Close()
		return nil, err
	}
	browser.Logger(opts.logger)
	browser.Trace(opts.trace)
	return browser, nil
}
