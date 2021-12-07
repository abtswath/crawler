package crawler

import (
	"crawler/browser"
	"crawler/config"
	"crawler/filter"
	"crawler/logger"
	"crawler/request"
	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
	"net/url"
	"sync"
	"time"
)

type Crawler struct {
	Browser *browser.Browser
	opts    *config.Option
	wg      sync.WaitGroup
	Result  []*request.Request
	lock    sync.Mutex
	Filter  filter.Filter
	timer   *time.Timer
	Logger  *logrus.Logger
	Pool    rod.PagePool
}

func NewCrawler(opts *config.Option) (*Crawler, error) {
	f := filter.NewDefaultFilter()
	f.RootHost = opts.Target.Host
	crawler := &Crawler{
		opts:   opts,
		wg:     sync.WaitGroup{},
		Filter: f,
		Logger: logger.New(),
		Pool:   rod.NewPagePool(opts.PoolSize),
	}

	var err error
	crawler.Browser, err = browser.NewBrowser(
		opts.BrowserPath,
		opts.Incognito,
		opts.Headless,
		opts.Proxy,
		opts.PoolSize,
		opts.PageTimeout,
	)
	crawler.Browser.Logger(crawler.Logger)
	if err != nil {
		return nil, err
	}

	return crawler, nil
}

func (c *Crawler) Run() {
	go func() {
		c.timer = time.AfterFunc(c.opts.Timeout, func() {
			c.Close()
		})
	}()
	defer c.Close()
	c.newJob(c.opts.Target)
	c.wg.Wait()
	var result []*request.Request
	for _, r := range c.Result {
		if !c.Filter.Exists(r) {
			result = append(result, r)
		}
	}
	c.Result = result
}

func (c *Crawler) newJob(target *url.URL) {
	// TODO. 页面递归嵌套问题
	c.wg.Add(1)
	go func() {
		c.Logger.Tracef("Start a new job: %s", target.String())
		defer c.wg.Done()
		page := c.Pool.Get(c.Browser.NewPage)
		defer c.Pool.Put(page)
		p := browser.NewPage(page, c.opts.Headers, target, browser.PageOption{
			IgnoreKeywords: c.opts.IgnoreKeywords,
			UploadFile:     c.opts.UploadFile,
			Logger:         c.Logger,
		})
		err := p.Run()
		if err != nil {
			c.Logger.Debugf("Page running error: %s", err)
			return
		}
		c.lock.Lock()
		defer c.lock.Unlock()
		c.Result = append(c.Result, p.Result...)
		for _, r := range p.Result {
			if c.Filter.Allow(r) &&
				!request.ShouldIgnoreRequest(*r, c.opts.IgnoreKeywords) &&
				!c.Filter.Exists(r) &&
				!c.Filter.Static(r) {
				c.newJob(r.URL)
			}
		}
	}()
}

func (c *Crawler) Close() error {
	if c.timer != nil {
		c.timer.Stop()
	}
	//c.Pool.Cleanup(func(page *rod.Page) {
	//	page.MustClose()
	//})
	return c.Browser.Close()
}
