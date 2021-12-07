package browser

import (
	"crawler/config"
	"crawler/filter"
	"crawler/request"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/sirupsen/logrus"
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

	Logger *logrus.Logger
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
		ctx.ContinueRequest(&proto.FetchContinueRequest{})
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
			p.collectURLFromResponse(ctx)
			break
		}
		p.addResult(r)
	})
}

func (p *Page) addResult(r *request.Request) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Result = append(p.Result, r)
}

func (p *Page) Run() error {
	_, err := p.Expose("collectURL", func(json gson.JSON) (interface{}, error) {
		r, err := request.NewRequestFromDOM(json.String(), p.MustInfo().URL)
		if err != nil {
			return nil, nil
		}
		p.addResult(r)
		return nil, nil
	})
	if err != nil {
		return err
	}
	err = p.hijack()
	if err != nil {
		return err
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
		p.opts.Logger.Debugf("Set extra headers error: %s", err)
		return err
	}
	defer cleanup()
	p.opts.Logger.Tracef("Ready for navigating %s", p.target.String())
	err = p.Navigate(p.target.String())
	if err != nil {
		p.opts.Logger.Tracef("Navigate error: %s", err)
		return err
	}
	err = p.WaitLoad()
	if err != nil {
		return err
	}

	p.wg.Add(1)
	p.collectURL()
	p.wg.Add(3)
	p.fillForm()
	p.wg.Wait()
	return nil
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
	go p.collectFromTagA()
}

func (p *Page) collectFromTagA() {
	defer p.wg.Done()
	elements, err := p.ElementsByJS(rod.Eval(`document.querySelectorAll('a[href]')`))
	if err != nil {
		p.opts.Logger.Debugf("Get tag a error: %s", err)
		return
	}
	pageInfo, err := p.Info()
	if err != nil {
		p.opts.Logger.Debugf("Get page info error: %s", err)
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
		r, err := request.NewRequestFromDOM(*href, pageInfo.URL)
		if err != nil {
			continue
		}
		p.addResult(r)
	}
}

func (p *Page) Close() error {
	if p.router != nil {
		_ = p.router.Stop()
	}
	return p.Page.Close()
}
