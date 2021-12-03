package crawler

import (
	"crawler/browser"
	"crawler/config"
	"crawler/filter"
	"crawler/utils"
	"net/url"
	"sync"
)

type Crawler struct {
	Browser *browser.Browser
	opts    *config.Options
	wg      sync.WaitGroup
	Result  []*Request
	lock    sync.Mutex
	Filter  filter.Filter
}

func NewCrawler(opts *config.Options) (*Crawler, error) {
	crawler := &Crawler{
		opts: opts,
		wg:   sync.WaitGroup{},
	}

	var err error
	crawler.Browser, err = browser.NewBrowser(
		opts.BrowserPath,
		opts.Incognito,
		opts.Headless,
		opts.Proxy,
		opts.PoolSize,
	)
	if err != nil {
		return nil, err
	}

	return crawler, nil
}

func (c *Crawler) Run() error {
	c.newJob(c.opts.Target)
	c.wg.Wait()
	return nil
}

func (c *Crawler) newJob(target *url.URL) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		page := c.Browser.Pool.Get(c.Browser.NewPage)
		defer c.Browser.Pool.Put(page)
		p := browser.NewPage(page, c.opts.Headers, target, browser.PageOption{
			Timeout:        5,
			IgnoreKeywords: c.opts.IgnoreKeywords,
		})
		p.Run()
		c.lock.Lock()
		defer c.lock.Unlock()
		c.Result = append(c.Result, p.Result...)
		for _, request := range p.Result {
			if !c.Filter.Exists(request) &&
				!c.Filter.Static(request) &&
				!utils.ShouldIgnoreRequest(*request, c.opts.IgnoreKeywords) {
				c.newJob(request.URL)
			}
		}
	}()
}

func (c *Crawler) Close() error {
	return c.Browser.Close()
}
