package crawler

import (
	"github.com/go-rod/rod"
	"sync"
)

type Crawler struct {
	browser *rod.Browser
	opts    Options
	pool    rod.PagePool
	wg      sync.WaitGroup
}

func NewCrawler(opts Options) (*Crawler, error) {
	crawler := &Crawler{
		opts: opts,
		pool: rod.NewPagePool(opts.PoolSize),
		wg:   sync.WaitGroup{},
	}
	err := crawler.startBrowser()
	if err != nil {
		return nil, err
	}
	return crawler, nil
}

func (c *Crawler) Run() {
	c.newJob(c.opts.Target.String())
	c.wg.Wait()
}

func (c *Crawler) Close() error {
	c.pool.Cleanup(func(page *rod.Page) {
		page.MustClose()
	})
	return c.browser.Close()
}
