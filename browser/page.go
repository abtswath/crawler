package browser

import (
	"crawler"
	"crawler/config"
	"crawler/utils"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

type PageOption struct {
	Timeout time.Duration

	IgnoreKeywords []string
}

type Page struct {
	*rod.Page
	Headers map[string]string
	router  *rod.HijackRouter
	wg      sync.WaitGroup
	Result  []*crawler.Request
	lock    sync.Mutex
	opts    PageOption
	target  *url.URL
}

func NewPage(page *rod.Page, headers map[string]string, target *url.URL, opts PageOption) *Page {
	page.Timeout(opts.Timeout)

	p := &Page{
		target:  target,
		opts:    opts,
		Headers: headers,
	}

	return p
}

func (p *Page) hijack() error {
	p.router = p.HijackRequests()
	return p.router.Add("", "", func(ctx *rod.Hijack) {
		request := crawler.NewRequestFromHijackRequest(ctx.Request, p.Headers)
		if utils.ShouldIgnoreRequest(*request, p.opts.IgnoreKeywords) {
			ctx.Skip = true
			ctx.Response.Fail(proto.NetworkErrorReasonAborted)
			p.addResult(request)
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
		p.addResult(request)
	})
}

func (p *Page) addResult(request *crawler.Request) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Result = append(p.Result, request)
}

func (p *Page) Run() {
	defer p.Close()
	err := p.hijack()
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

			p.addResult(crawler.NewRequestFromHijackRequest(ctx.Request, p.Headers))
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
		href, err := element.Property("href")
		if err != nil {
			continue
		}
		request, err := crawler.NewRequestFromDOM(href.String(), pageInfo.URL)
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
