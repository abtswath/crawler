package browser

import (
	"crawler/config"
	"crawler/filter"
	"crawler/request"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/ysmood/gson"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

type PageOption struct {
	IgnoreKeywords []string

	UploadFile string

	Filter filter.Filter
}

type Page struct {
	*rod.Page
	Headers map[string]string
	router  *rod.HijackRouter
	wg      sync.WaitGroup
	Result  []*request.Request
	lock    sync.Mutex
	opts    PageOption
	target  *url.URL
}

func NewPage(page *rod.Page, headers map[string]string, target *url.URL, opts PageOption) *Page {
	p := &Page{
		Page:    page,
		target:  target,
		opts:    opts,
		Headers: headers,
	}

	return p
}

func (p *Page) hijack() error {
	p.router = p.HijackRequests()
	return p.router.Add("", "", func(ctx *rod.Hijack) {
		r := request.NewRequestFromHijackRequest(ctx.Request, p.Headers)
		if request.ShouldIgnoreRequest(*r, p.opts.IgnoreKeywords) {
			ctx.Skip = true
			ctx.Response.Fail(proto.NetworkErrorReasonAborted)
			p.addResult(r)
			return
		}

		switch ctx.Request.Type() {
		case proto.NetworkResourceTypeImage:
			fallthrough
		case proto.NetworkResourceTypeMedia:
			fallthrough
		case proto.NetworkResourceTypeFont:
			fallthrough
		case proto.NetworkResourceTypeTextTrack:
			fallthrough
		case proto.NetworkResourceTypeSignedExchange:
			fallthrough
		case proto.NetworkResourceTypeCSPViolationReport:
			fallthrough
		case proto.NetworkResourceTypePing:
			fallthrough
		case proto.NetworkResourceTypePreflight:
			fallthrough
		case proto.NetworkResourceTypeOther:
			ctx.Skip = true
			ctx.Response.Fail(proto.NetworkErrorReasonAborted)
			break
		case proto.NetworkResourceTypeDocument:
			fallthrough
		case proto.NetworkResourceTypeStylesheet:
			fallthrough
		case proto.NetworkResourceTypeScript:
			fallthrough
		case proto.NetworkResourceTypeXHR:
			fallthrough
		case proto.NetworkResourceTypeFetch:
			fallthrough
		case proto.NetworkResourceTypeEventSource:
			fallthrough
		case proto.NetworkResourceTypeWebSocket:
			fallthrough
		case proto.NetworkResourceTypeManifest:
			ctx.ContinueRequest(&proto.FetchContinueRequest{})
			p.collectURLFromResponse(ctx)
			break
		}
		p.addResult(r)
	})
}

func (p *Page) addResult(request *request.Request) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Result = append(p.Result, request)
}

func (p *Page) Run() {
	defer p.Close()
	_, err := p.Expose("collectURL", func(json gson.JSON) (interface{}, error) {
		request, err := request.NewRequestFromDOM(json.String(), p.MustInfo().URL)
		if err != nil {
			return nil, nil
		}
		p.addResult(request)
		return nil, nil
	})
	if err != nil {
		return
	}
	err = p.Navigate(p.target.String())
	if err != nil {
		return
	}
	go p.EachEvent(func(e *proto.PageFrameRequestedNavigation) {
		_ = p.StopLoading()
	})()
	err = p.hijack()
	if err != nil {
		return
	}
	go p.router.Run()

	var headers []string
	for key, value := range p.Headers {
		if key != "Host" {
			headers = append(headers, key, value)
		}
	}
	cleanup, err := p.SetExtraHeaders(headers)
	if err != nil {
		return
	}
	defer cleanup()
	err = p.WaitLoad()
	if err != nil {
		return
	}

	p.collectURL()
	p.fillForm()
	p.wg.Wait()
}

func (p *Page) collectURLFromResponse(ctx *rod.Hijack) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		body := ctx.Response.Body()
		regex := regexp.MustCompile(config.SuspectURLRegex)
		result := regex.FindAllString(body, -1)
		for _, u := range result {
			u = u[1 : len(u)-1]
			urlLowerCase := strings.ToLower(u)
			if strings.HasPrefix(urlLowerCase, "image/x-icon") || strings.HasPrefix(urlLowerCase, "text/css") || strings.HasPrefix(urlLowerCase, "text/javascript") {
				continue
			}

			p.addResult(request.NewRequestFromHijackRequest(ctx.Request, p.Headers))
		}
	}()
}

func (p *Page) collectURL() {
	p.wg.Add(1)
	go p.collectFromTagA()
}

func (p *Page) collectFromTagA() {
	defer p.wg.Done()
	elements, err := p.ElementsByJS(rod.Eval(`document.querySelectorAll('a[href]')`))
	if err != nil {
		return
	}
	pageInfo, err := p.Info()
	if err != nil {
		return
	}
	for _, element := range elements {
		href, err := element.Attribute("href")
		if err != nil || href == nil {
			continue
		}
		if strings.HasPrefix(*href, "javascript:") {
			continue
		}
		request, err := request.NewRequestFromDOM(*href, pageInfo.URL)
		if err != nil {
			continue
		}
		p.addResult(request)
	}
}

func (p *Page) Close() error {
	if p.router != nil {
		_ = p.router.Stop()
	}
	return p.Page.Close()
}
