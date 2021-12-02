package crawler

import (
	"crawler/utils"
	"github.com/go-rod/rod"
	"net/url"
)

type Request struct {
	URL     *url.URL
	Method  string
	Headers map[string]string
	Body    string
}

func NewRequestFromHijackRequest(request *rod.HijackRequest, extraHeaders map[string]string) *Request {
	headers := map[string]string{}
	for key, value := range request.Headers() {
		headers[key] = value.String()
	}
	for key, value := range extraHeaders {
		headers[key] = value
	}

	return &Request{
		URL:     request.URL(),
		Method:  request.Method(),
		Headers: headers,
		Body:    request.Body(),
	}
}

func (r Request) UniqueID() string {
	return utils.Hash(r.Method + r.URL.String() + r.Body)
}
