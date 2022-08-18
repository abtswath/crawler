package browser

import (
	"log"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type Browser struct {
	browser *rod.Browser
	opts    Option
	pool    rod.PagePool
	logger  *log.Logger
}

func New(opts Option) *Browser {
	return &Browser{
		opts:   opts,
		logger: opts.Logger,
	}
}

func (b *Browser) controlURL() (string, error) {
	l := launcher.New().
		Leakless(true).
		Headless(b.opts.Headless).
		NoSandbox(true).
		Set("disable-gpu").
		Set("disable-web-security").
		Set("disable-xss-auditor").
		Set("disable-setuid-sandbox").
		Set("allow-running-insecure-content").
		Set("disable-popup-blocking").
		Set("disable-webgl").
		Set("ignore-certificate-errors").
		Set("disable-popup-blocking").
		Set("disable-images").
		Set("incognito")
	if b.opts.Bin != "" {
		l.Bin(b.opts.Bin)
	}
	if b.opts.Proxy != "" {
		l.Proxy(b.opts.Proxy)
	}
	return l.Launch()
}

func (b *Browser) Start() error {
	controlURL, err := b.controlURL()
	if err != nil {
		return err
	}
	b.browser = rod.New().
		Trace(true).
		Logger(b.logger).
		ControlURL(controlURL)
	if err := b.browser.Connect(); err != nil {
		b.browser.Close()
		return err
	}
	if err := b.browser.SetCookies(b.opts.Cookies); err != nil {
		b.browser.Close()
		return err
	}
	return nil
}

func (b *Browser) NewPage() *rod.Page {
	b.pool = rod.NewPagePool(b.opts.MaxPageCount)
	page := b.pool.Get(b.createPage)
	defer b.pool.Put(page)
	_ = page.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: b.opts.UserAgent})
	_, _ = page.SetExtraHeaders(b.opts.Headers)

	go page.EachEvent(func(e *proto.PageJavascriptDialogOpening) {
		_ = proto.PageHandleJavaScriptDialog{
			Accept: true,
		}.Call(page)
	})()
	return page
}

func (b *Browser) createPage() *rod.Page {
	return b.browser.MustPage()
}

func (b *Browser) Done() <-chan struct{} {
	return b.browser.GetContext().Done()
}

func (b *Browser) cleanup(page *rod.Page) {
	_ = page.Close()
}

func (b *Browser) Close() error {
	b.pool.Cleanup(b.cleanup)
	return b.browser.Close()
}
