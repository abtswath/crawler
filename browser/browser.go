package browser

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"time"
)

type Browser struct {
	*rod.Browser

	Pool rod.PagePool

	PageTimeout time.Duration
}

func NewBrowser(bin string, incognito bool, headless bool, proxy string, poolSize int, pageTimeout time.Duration) (*Browser, error) {
	l := launcher.New().
		Leakless(true).
		Headless(headless).
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
	if incognito {
		l.Set("incognito")
	}
	if proxy != "" {
		l.Proxy(proxy)
	}
	if bin != "" {
		l.Bin(bin)
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
	b := &Browser{
		Browser:     browser,
		Pool:        rod.NewPagePool(poolSize),
		PageTimeout: pageTimeout,
	}

	// TODO. Is the HandleAuth correctly?
	go b.MustHandleAuth("username", "password")()

	return b, nil
}

func (b *Browser) NewPage() *rod.Page {
	page, _ := b.Page(proto.TargetCreateTarget{
		Width:  1920,
		Height: 1080,
	})
	page.Timeout(b.PageTimeout)
	_, _ = proto.PageAddScriptToEvaluateOnNewDocument{
		Source: injectionScript,
	}.Call(page)
	return page
}

func (b *Browser) Close() error {
	b.Pool.Cleanup(func(page *rod.Page) {
		_ = page.Close()
	})
	return b.Browser.Close()
}
