package crawler

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func (c *Crawler) newPage() *rod.Page {
	page, err := c.browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil
	}
	go page.EachEvent(func(e *proto.TargetTargetCreated) {
		page.MustEval(injectionScript)
	}, func(e *proto.PageFrameRequestedNavigation) {
		// TODO. Collect URL
		//navigationURL, _ := url.Parse(e.URL)
		//b.queue <- &URL{
		//	URL:    navigationURL,
		//	Method: "GET",
		//}
		_ = page.StopLoading()
	})()
	return page
}

var injectionScript = `
document.body.addEventListener('click', function () {
	var target = event.target;
	if (target.nodeName.toLocaleLowerCase() === 'a') {
		event.preventDefault();
	}
});
`
