package crawler

import (
	"context"
	"crawler/collection"
	"crawler/filter"
	"crawler/logger"
	"crawler/request"
	"crawler/tab"
	"errors"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/sirupsen/logrus"
	"net/url"
	"sync"
	"time"
)

// TODO. Disable select file window and javascript injection moment

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
	Filter         filter.Filter
	BrowserTrace   bool
	Cookies        []*proto.NetworkCookieParam
	LogLevel       logrus.Level
}

type Crawler struct {
	browser    *rod.Browser
	opts       Option
	wg         sync.WaitGroup
	Result     *collection.Collection
	Filter     filter.Filter
	Logger     *logrus.Logger
	Pool       rod.PagePool
	context    context.Context
	cancelFunc context.CancelFunc
}

const (
	DefaultTimeout     = time.Minute
	DefaultPageTimeout = time.Second * 3
)

func New(opts Option) (*Crawler, error) {
	if opts.Target == nil {
		return nil, errors.New("invalid target")
	}
	if opts.Filter == nil {
		opts.Filter = filter.NewDefaultFilter(opts.Target.Host)
	}
	if opts.Timeout == 0 {
		opts.Timeout = DefaultTimeout
	}
	if opts.PageTimeout == 0 {
		opts.PageTimeout = DefaultPageTimeout
	}
	if opts.Headers == nil {
		opts.Headers = map[string]string{}
	}
	if opts.PoolSize == 0 {
		opts.PoolSize = 15
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), opts.Timeout)

	crawler := &Crawler{
		opts:       opts,
		wg:         sync.WaitGroup{},
		Filter:     opts.Filter,
		Logger:     logger.New(opts.LogLevel),
		Pool:       rod.NewPagePool(opts.PoolSize),
		context:    ctx,
		cancelFunc: cancelFunc,
		Result:     collection.NewCollection(opts.Target.Host),
	}

	var err error
	crawler.browser, err = newBrowser(browserOption{
		bin:         opts.BrowserPath,
		incognito:   opts.Incognito,
		headless:    opts.Headless,
		proxy:       opts.Proxy,
		pageTimeout: opts.PageTimeout,
		logger:      crawler.Logger,
		trace:       opts.BrowserTrace,
		cookies:     opts.Cookies,
	})
	if err != nil {
		return nil, err
	}

	return crawler, nil
}

func (c *Crawler) page() *rod.Page {
	page, err := c.browser.Page(proto.TargetCreateTarget{
		URL:    "",
		Width:  1920,
		Height: 1080,
	})
	if err != nil {
		c.Logger.Traceln("Create page error: %s", err)
		return nil
	}
	_ = proto.PageSetInterceptFileChooserDialog{Enabled: true}.Call(page)
	go page.EachEvent(func(e *proto.PageFrameRequestedNavigation) {
		_ = page.StopLoading()
	}, func(e *proto.PageJavascriptDialogOpening) {
		_ = proto.PageHandleJavaScriptDialog{
			Accept:     true,
			PromptText: "",
		}.Call(page)
	})()
	return page
}

func (c *Crawler) newJob(target *url.URL) {
	defer c.wg.Done()
	p := c.Pool.Get(c.page)
	defer c.Pool.Put(p)
	t := tab.New(p, target, tab.Option{
		IgnoreKeywords: c.opts.IgnoreKeywords,
		UploadFile:     c.opts.UploadFile,
		Logger:         c.Logger,
		Headers:        c.opts.Headers,
		Timeout:        c.opts.PageTimeout,
	})
	err := t.Run()
	if err != nil {
		c.Logger.Debugf("Tab running error: %s", err)
		return
	}
	c.Result.Put(t.Result...)
	for _, r := range t.Result {
		if c.Filter.Allow(r) &&
			!c.Filter.Exists(r) &&
			!c.Filter.Static(r) &&
			!request.ShouldIgnoreRequest(*r, c.opts.IgnoreKeywords) {
			c.wg.Add(1)
			go c.newJob(r.URL)
		}
	}
}

func (c *Crawler) run() {
	defer c.cancelFunc()
	c.wg.Add(1)
	go c.newJob(c.opts.Target)
	c.wg.Wait()
}

func (c *Crawler) Run() {
	go c.run()
	select {
	case <-c.context.Done():
		c.Logger.Traceln("Done...")
		//c.browser.Close()
		return
	}
}

func (c *Crawler) close() error {
	return c.browser.Close()
}

func (c *Crawler) closePage(page *rod.Page) {
	_ = page.Close()
}

func (c *Crawler) Close() error {
	c.Pool.Cleanup(c.closePage)
	return c.close()
}
