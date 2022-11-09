package crawler

import (
	"context"
	"crawler/pkg/browser"
	"crawler/pkg/constants"
	"crawler/pkg/filter"
	"crawler/pkg/model"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/go-rod/rod/lib/proto"
)

type Crawler struct {
	target      url.URL
	browser     *browser.Browser
	opts        Option
	PageTimeout time.Duration
	logger      *log.Logger
	trees       model.Trees
	ctx         context.Context
	cancel      context.CancelFunc
	lock        sync.Mutex
	filter      filter.Filter
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

	c := Crawler{
		target:      target,
		browser:     b,
		PageTimeout: opts.PageTimeout,
		logger:      opts.Logger,
		trees:       model.Trees{},
		ctx:         ctx,
		cancel:      cancel,
	}
	c.filter = filter.New(c.target.Host, c.opts.Exclusions)
	return &c, nil
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

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			case req, ok := <-page.Request():
				if ok {
					root := c.treeNode(req.Method)
					if root.Get(req.URL.Path) != nil {
						return
					}
					root.Put(req.URL.Path, req.URL, req.ResourceType)
					if c.filter.Can(req) {
						go c.newTask(req.URL.String())
					}
				}
			}
		}
	}()

	if err := page.Navigate(address); err != nil {
		return
	}

	page.Collect()
}

func (c *Crawler) treeNode(method string) *model.Node {
	if node := c.trees.Get(method); node != nil {
		return node
	}
	root := &model.Node{
		URL: url.URL{
			Scheme: c.target.Scheme,
			Host:   c.target.Host,
			Path:   "/",
		},
		ResourceType: proto.NetworkResourceTypeDocument,
	}
	c.lock.Lock()
	c.trees = append(c.trees, model.NewTree(method, root))
	c.lock.Unlock()
	return root
}

func (c *Crawler) Get() []model.Result {
	result := []model.Result{}
	for _, tree := range c.trees {
		result = append(result, collectResult(tree.Root, tree.Method)...)
	}
	return result
}

func (c *Crawler) Close() error {
	if c.browser != nil {
		return c.browser.Close()
	}
	return nil
}

func collectResult(node *model.Node, method string) []model.Result {
	result := []model.Result{}
	if len(node.Children()) <= 0 {
		result = append(result, model.Result{
			URL:          node.URL.String(),
			Method:       method,
			ResourceType: node.ResourceType,
		})
	}
	for _, n := range node.Children() {
		result = append(result, collectResult(n, method)...)
	}
	return result
}
