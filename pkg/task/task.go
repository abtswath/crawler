package task

import (
	"context"
	"crawler/pkg/browser"
	"crawler/pkg/constants"
	"log"
	"net/url"
	"time"
)

type Task struct {
	target      string
	browser     *browser.Browser
	opts        Option
	PageTimeout time.Duration
	logger      *log.Logger
	Collection  Collection
	ctx         context.Context
	cancel      context.CancelFunc
}

func New(target string, opts Option) (*Task, error) {

	t, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

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

	c := Collection{}
	c.filter = &BasicFilter{
		domains:    []string{t.Hostname()},
		collection: &c,
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)

	return &Task{
		target:      target,
		browser:     b,
		PageTimeout: opts.PageTimeout,
		logger:      opts.Logger,
		Collection:  c,
		ctx:         ctx,
		cancel:      cancel,
	}, nil
}

func (t *Task) Run() {
	defer t.Close()
	go t.newTask(t.target)
	for {
		select {
		case <-t.ctx.Done():
			return
		}
	}
}

func (t *Task) newTask(url string) {
	page := browser.NewPage(t.browser, browser.PageOption{
		Timeout: t.opts.PageTimeout,
	})

	if err := page.Navigate(url); err != nil {
		return
	}

	go func() {
		for {
			select {
			case <-t.ctx.Done():
				return
			case req, ok := <-page.Request():
				if ok {
					if !t.Collection.Has(req) {
						t.Collection.Put(req)
						if t.Collection.filter.Can(req) {
							go t.newTask(t.target)
						}
					}
				}
			}
		}
	}()
	page.Collect()
}

func (t *Task) Close() error {
	return t.browser.Close()
}
