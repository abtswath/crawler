package crawler

import (
	"crawler/filter"
	"crawler/logger"
	"crawler/request"
	"crawler/tab"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/sirupsen/logrus"
	"net/url"
	"sync"
	"time"
)

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
}

var s = []string{"http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003", "http://localhost:9003"}
var index = 0

type Crawler struct {
	browser *rod.Browser
	opts    Option
	wg      sync.WaitGroup
	Result  []*request.Request
	lock    sync.Mutex
	Filter  filter.Filter
	timer   *time.Timer
	Logger  *logrus.Logger
	Pool    rod.PagePool
}

func New(opts Option) (*Crawler, error) {
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
	crawler.browser, err = newBrowser(browserOption{
		bin:         opts.BrowserPath,
		incognito:   opts.Incognito,
		headless:    opts.Headless,
		proxy:       opts.Proxy,
		pageTimeout: opts.PageTimeout,
		logger:      crawler.Logger,
	})
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
	c.wg.Add(1)
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

func (c *Crawler) page() *rod.Page {
	page, _ := c.browser.Page(proto.TargetCreateTarget{
		URL:    "",
		Width:  1920,
		Height: 1080,
	})
	//page.Timeout(c.opts.PageTimeout)
	//_, _ = proto.PageAddScriptToEvaluateOnNewDocument{
	//	Source: injectionScript,
	//}.Call(page)
	//go page.EachEvent(func(e *proto.PageFrameRequestedNavigation) {
	//	_ = page.StopLoading()
	//})()
	return page
}

func (c *Crawler) newJob(target *url.URL) {
	c.Logger.Tracef("Start a new job: %s", target.String())
	defer c.wg.Done()
	p := c.Pool.Get(c.page)
	defer c.Pool.Put(p)
	t := tab.New(p, target, tab.Option{
		IgnoreKeywords: c.opts.IgnoreKeywords,
		UploadFile:     c.opts.UploadFile,
		Logger:         c.Logger,
		Headers:        c.opts.Headers,
	})
	err := t.Run()
	if err != nil {
		c.Logger.Debugf("Page running error: %s", err)
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Result = append(c.Result, t.Result...)
	for _, r := range t.Result {
		if c.Filter.Allow(r) &&
			!request.ShouldIgnoreRequest(*r, c.opts.IgnoreKeywords) &&
			!c.Filter.Exists(r) &&
			!c.Filter.Static(r) {
			c.wg.Add(1)
			go c.newJob(r.URL)
		}
	}
}

func (c *Crawler) Close() error {
	if c.timer != nil {
		c.timer.Stop()
	}
	fmt.Println(len(c.Pool))
	c.Pool.Cleanup(func(p *rod.Page) {
		p.MustClose()
	})
	return c.browser.Close()
}
