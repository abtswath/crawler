package tab

import (
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

type Option struct {
	IgnoreKeywords []string
	UploadFile     string
	Filter         filter.Filter
	Logger         *logrus.Logger
	Headers        map[string]string
}

type Tab struct {
	Page           *rod.Page
	Headers        map[string]string
	router         *rod.HijackRouter
	wg             sync.WaitGroup
	Result         []*request.Request
	lock           sync.Mutex
	target         *url.URL
	uploadFile     string
	ignoreKeywords []string
	logger         *logrus.Logger
	filter         filter.Filter
}

func New(page *rod.Page, target *url.URL, opts Option) *Tab {
	p := &Tab{
		Page:       page,
		target:     target,
		logger:     opts.Logger,
		uploadFile: opts.UploadFile,
		filter:     opts.Filter,
		Headers:    opts.Headers,
	}

	return p
}

func (t *Tab) hijack() error {
	t.router = t.Page.HijackRequests()
	return t.router.Add("", "", func(ctx *rod.Hijack) {
		r := request.NewRequestFromHijackRequest(ctx.Request, t.Headers)
		if request.ShouldIgnoreRequest(*r, t.ignoreKeywords) {
			ctx.Skip = true
			ctx.Response.Fail(proto.NetworkErrorReasonAborted)
			t.addResult(r)
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
			t.collectURLFromResponse(ctx)
			break
		}
		ctx.ContinueRequest(&proto.FetchContinueRequest{})
		t.addResult(r)
	})
}

func (t *Tab) addResult(r *request.Request) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.Result = append(t.Result, r)
}

func (t *Tab) Run() error {
	_, err := t.Page.Expose("collectURL", func(json gson.JSON) (interface{}, error) {
		r, err := request.NewRequestFromDOM(json.String(), t.Page.MustInfo().URL)
		if err != nil {
			return nil, nil
		}
		t.addResult(r)
		return nil, nil
	})
	if err != nil {
		return err
	}
	err = t.hijack()
	if err != nil {
		return err
	}
	go t.router.Run()

	var headers []string
	for key, value := range t.Headers {
		if key != "Host" {
			headers = append(headers, key, value)
		}
	}
	cleanup, err := t.Page.SetExtraHeaders(headers)
	if err != nil {
		t.logger.Debugf("Set extra headers error: %s", err)
		return err
	}
	defer cleanup()
	t.logger.Tracef("Ready for navigating %s", t.target.String())
	err = t.Page.Navigate(t.target.String())
	if err != nil {
		t.logger.Tracef("Navigate error: %s", err)
		return err
	}
	err = t.Page.WaitLoad()
	if err != nil {
		return err
	}

	t.wg.Add(1)
	t.collectURL()
	t.wg.Add(3)
	t.fillForm()
	t.wg.Wait()
	return nil
}

func (t *Tab) collectURLFromResponse(ctx *rod.Hijack) {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		body := ctx.Response.Body()
		regex := regexp.MustCompile(SuspectURLRegex)
		result := regex.FindAllString(body, -1)
		for _, u := range result {
			u = u[1 : len(u)-1]
			urlLowerCase := strings.ToLower(u)
			if strings.HasPrefix(urlLowerCase, "image/x-icon") || strings.HasPrefix(urlLowerCase, "text/css") || strings.HasPrefix(urlLowerCase, "text/javascript") {
				continue
			}

			t.addResult(request.NewRequestFromHijackRequest(ctx.Request, t.Headers))
		}
	}()
}

func (t *Tab) collectURL() {
	go t.collectFromTagA()
}

func (t *Tab) collectFromTagA() {
	defer t.wg.Done()
	elements, err := t.Page.ElementsByJS(rod.Eval(`document.querySelectorAll('a[href]')`))
	if err != nil {
		t.logger.Debugf("Get tag a error: %s", err)
		return
	}
	pageInfo := t.Page.MustInfo()
	if err != nil {
		t.logger.Debugf("Get page info error: %s", err)
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
		t.addResult(r)
	}
}
