package crawler

import (
	"context"
	"crawler/pkg/browser"
	"crawler/pkg/constants"
	"crawler/pkg/model"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Crawler struct {
	target      url.URL
	browser     *browser.Browser
	opts        Option
	PageTimeout time.Duration
	logger      *log.Logger
	trees       trees
	ctx         context.Context
	cancel      context.CancelFunc
}

func New(target url.URL, opts Option) (*Crawler, error) {

	for _, fn := range []SetOptionValueFunc{
		WithTimeout(constants.DefaultTimeout),
		WithHeadless(constants.DefaultHeadless),
		WithMaxPageCount(constants.DefaultMaxPageCount),
		WithPageTimeout(constants.DefaultPageTimeout),
		WithUserAgent(constants.DefaultUserAgent),
		WithLogger(log.Default()),
	} {
		fn(&opts)
	}

	b := browser.New(browser.Option{
		MaxPageCount: opts.MaxPageCount,
		Logger:       opts.Logger,
		Headless:     opts.Headless,
		Bin:          opts.Bin,
		Proxy:        opts.Proxy,
		UserAgent:    opts.UserAgent,
		Cookies:      opts.Cookies,
		Headers:      opts.Headers,
	})
	if err := b.Start(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)

	return &Crawler{
		target:      target,
		browser:     b,
		PageTimeout: opts.PageTimeout,
		logger:      opts.Logger,
		trees:       trees{},
		ctx:         ctx,
		cancel:      cancel,
	}, nil
}

func (c *Crawler) Run() {
	defer c.Close()
	go c.newTask(c.target.String())
	<-c.ctx.Done()
}

func (c *Crawler) newTask(address string) {
	page := browser.NewPage(c.browser, browser.PageOption{
		Timeout: c.opts.PageTimeout,
	})

	if err := page.Navigate(address); err != nil {
		return
	}

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			case req, ok := <-page.Request():
				if ok {
					root := c.treeNode(req.Method)
					if !root.has(req.URL.Path) {
						root.put(req.URL.Path, req)
						if c.can(req) {
							go c.newTask(req.URL.String())
						}
					}
				}
			}
		}
	}()
	page.Collect()
}

func (c *Crawler) treeNode(method string) *node {
	root := c.trees.get(method)
	if root == nil {
		root = new(node)
		root.request = model.Request{
			URL: url.URL{
				Scheme: c.target.Scheme,
				Host:   c.target.Host,
				Path:   "/",
			},
			Method:       http.MethodGet,
			ResourceType: constants.ResourceTypeDocument,
		}
		c.trees = append(c.trees, tree{method: method, root: root})
	}
	return root
}

func (c *Crawler) can(request model.Request) bool {
	if request.URL.Host != c.target.Host {
		return false
	}
	for _, kw := range c.opts.Exclusions {
		if strings.Contains(request.URL.Path, kw) {
			return false
		}
	}
	return true
}

func (c *Crawler) Get() []model.Request {
	result := []model.Request{}
	for _, tree := range c.trees {
		result = append(result, tree.root.all()...)
	}
	return result
}

func (c *Crawler) Close() error {
	return c.browser.Close()
}
