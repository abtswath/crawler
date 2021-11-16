package crawler

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"net/url"
)

func (c *Crawler) launcher() (string, error) {
	return launcher.New().
		Leakless(true).
		Set("incognito").
		Set("ignore-certificate-errors").
		Set("disable-webgl").
		Set("disable-popup-blocking").
		Set("disable-images").
		Set("user-agent", c.opts.UserAgent).
		Launch()
}

func (c *Crawler) startBrowser() error {
	u, err := c.launcher()
	if err != nil {
		return err
	}
	c.browser = rod.New().ControlURL(u)
	err = c.browser.Connect()
	if err != nil {
		c.browser.MustClose()
		return err
	}
	err = c.hijack()
	if err != nil {
		c.browser.MustClose()
		return err
	}
	return nil
}

func (c *Crawler) hijack() error {
	router := c.browser.HijackRequests()
	patternURL := &url.URL{
		Host:   c.opts.Target.Host,
		Scheme: c.opts.Target.Scheme,
		Path:   "*",
	}
	return router.Add(patternURL.String(), "", func(hijack *rod.Hijack) {
		hijack.ContinueRequest(&proto.FetchContinueRequest{})
		// TODO. Collect URL
		//u := &URL{
		//	URL:    hijack.Request.URL(),
		//	Method: hijack.Request.Method(),
		//}
		//if !b.Collection.Exists(u) {
		//	b.Collection.Put(u)
		//	if strings.ToLower(hijack.Response.Headers().Get("content-type")) == "text/html" {
		//		b.queue <- u
		//	}
		//}
	})
}
