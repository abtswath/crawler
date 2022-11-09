package browser

import (
	"net/http"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func (p *Page) collectURLFromForms() {
	defer p.wg.Done()
	forms, err := p.page.ElementsByJS(rod.Eval(`document.querySelectorAll('form[action]')`))
	if err != nil {
		return
	}
	for _, form := range forms {
		action, err := form.Attribute("action")
		if err != nil {
			continue
		}
		method, err := form.Attribute("method")
		if err != nil || method == nil {
			p.send(*action, http.MethodGet, proto.NetworkResourceTypeDocument)
		} else {
			p.send(*action, *method, proto.NetworkResourceTypeDocument)
		}
	}
}

func (p *Page) collectURLFromObject() {
	defer p.wg.Done()
	objects, err := p.page.ElementsByJS(rod.Eval(`() => document.querySelectorAll('object[data]')`))
	if err != nil {
		return
	}
	for _, item := range objects {
		data, err := item.Attribute("data")
		if err != nil {
			continue
		}
		p.send(*data, http.MethodGet, proto.NetworkResourceTypeOther)
	}
}

func (p *Page) collectURLFromAnchors() {
	defer p.wg.Done()
	anchors, err := p.page.ElementsByJS(rod.Eval(`() => document.querySelectorAll('a[href]')`))
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
		p.send(*href, http.MethodGet, proto.NetworkResourceTypeDocument)
	}
}
