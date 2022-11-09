package browser

import (
	"crawler/pkg/model"
	"crawler/pkg/utils"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/ysmood/gson"
)

type Page struct {
	page    *rod.Page
	timeout time.Duration
	router  *rod.HijackRouter
	channel chan model.Request
	info    *proto.TargetTargetInfo
	wg      sync.WaitGroup
}

func NewPage(browser *Browser, opts PageOption) Page {
	return Page{
		page:    browser.NewPage(),
		timeout: opts.Timeout,
		channel: make(chan model.Request),
	}
}

func (p *Page) Navigate(url string) error {
	p.page.Timeout(p.timeout)
	wait := p.page.WaitNavigation(proto.PageLifecycleEventNameFirstMeaningfulPaint)
	p.handleEvent()
	p.injectScript()
	if err := p.page.Navigate(url); err != nil {
		return err
	}
	var err error
	p.info, err = p.page.Info()
	if err != nil {
		return err
	}
	if err := p.hijack(); err != nil {
		return err
	}
	wait()
	return nil
}

func (p *Page) Collect() {
	defer p.router.Stop()
	p.wg.Add(3)
	go p.collectURLFromForms()
	go p.collectURLFromAnchors()
	go p.collectURLFromObject()
	p.wg.Wait()
}

func (p *Page) handleEvent() {
	go p.page.EachEvent(func(e *proto.PageJavascriptDialogOpening) {
		_ = proto.PageHandleJavaScriptDialog{
			Accept: true,
		}.Call(p.page)
	}, func(e *proto.ConsoleMessageAdded) {
		fmt.Printf("%s: %s\n", e.Message.Level, e.Message.Text)
	})()
}

func (p *Page) hijack() error {
	p.router = p.page.HijackRequests()
	if err := p.router.Add("*", "", p.hijackHandler); err != nil {
		return err
	}
	go p.router.Run()
	return nil
}

func (p *Page) hijackHandler(ctx *rod.Hijack) {
	p.send(ctx.Request.URL().String(), ctx.Request.Method(), ctx.Request.Type())
	if ctx.Request.IsNavigation() {
		ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
		return
	}
	switch ctx.Request.Type() {
	case proto.NetworkResourceTypeFetch,
		proto.NetworkResourceTypeXHR,
		proto.NetworkResourceTypeWebSocket:
		ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
	}
	ctx.ContinueRequest(&proto.FetchContinueRequest{})
}

func (p *Page) Request() chan model.Request {
	return p.channel
}

func (p *Page) injectScript() {
	p.page.Eval(injectionScript)
	p.page.Expose("collectLink", func(g gson.JSON) (interface{}, error) {
		p.send(g.Str(), http.MethodGet, proto.NetworkResourceTypeDocument)
		return nil, nil
	})
}

func (p *Page) send(value string, method string, resourceType proto.NetworkResourceType) {
	u, err := utils.ParseURL(value, p.info.URL)
	if err != nil {
		return
	}
	p.channel <- model.Request{
		URL:          *u,
		Method:       method,
		ResourceType: resourceType,
	}
}
