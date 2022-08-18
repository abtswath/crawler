package browser

import (
	"crawler/pkg/constants"
	"crawler/pkg/model"
	"crawler/pkg/utils"
	"net/http"

	"github.com/go-rod/rod"
)

func (p *Page) collectURLFromAnchors() {
	p.wg.Done()
	anchors, err := p.page.ElementsByJS(rod.Eval(`() => document.getElementsByTagName('a')`))
	if err != nil {
		return
	}
	for _, anchor := range anchors {
		href, err := anchor.Attribute("href")
		if err != nil {
			continue
		}
		if href == nil {
			continue
		}
		u, err := utils.ParseURL(*href, p.info.URL)
		if err != nil {
			continue
		}
		p.channel <- model.Request{
			URL:          *u,
			Method:       http.MethodGet,
			ResourceType: constants.ResourceTypeDocument,
		}
	}
}
