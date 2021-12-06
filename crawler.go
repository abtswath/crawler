package crawler

import (
	"crawler/browser"
	"crawler/config"
	"crawler/filter"
	"crawler/request"
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
}

func NewCrawler(opts *config.Option) (*Crawler, error) {
	crawler := &Crawler{
		opts:   opts,
		wg:     sync.WaitGroup{},
		Filter: filter.NewDefaultFilter(),
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
	for _, request := range c.Result {
		if !c.Filter.Exists(request) {
			result = append(result, request)
		}
	}
	c.Result = result
}

func (c *Crawler) newJob(target *url.URL) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		page := c.Browser.Pool.Get(c.Browser.NewPage)
		defer c.Browser.Pool.Put(page)
		p := browser.NewPage(page, c.opts.Headers, target, browser.PageOption{
			IgnoreKeywords: c.opts.IgnoreKeywords,
			UploadFile:     c.opts.UploadFile,
		})
		p.Run()
		c.lock.Lock()
		defer c.lock.Unlock()
		c.Result = append(c.Result, p.Result...)
		for _, r := range p.Result {
			if !c.Filter.Exists(r) &&
				!request.ShouldIgnoreRequest(*r, c.opts.IgnoreKeywords) {
				c.newJob(r.URL)
			}
		}
	}()
}

func (c *Crawler) Close() error {
	if c.timer != nil {
		c.timer.Stop()
	}
	return c.Browser.Close()
}
