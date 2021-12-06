package crawler

import (
	"crawler/utils"
	"github.com/go-rod/rod"
	"net/url"
)

const (
	MethodGET     = "GET"
	MethodHEAD    = "HEAD"
	MethodPOST    = "POST"
	MethodPUT     = "PUT"
	MethodDELETE  = "DELETE"
	MethodOPTION  = "OPTION"
	MethodCONNECT = "CONNECT"
	MethodTRACE   = "TRACE"
	MethodPATCH   = "PATCH"
)

const (
	RequestTypeNormal = 1
	RequestTypeDOM
	RequestTypeComment
)

type Request struct {
	URL      *url.URL
	Method   string
	Headers  map[string]string
	Body     string
	UniqueID string
	Type     int
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
		URL:      request.URL(),
		Method:   request.Method(),
		Headers:  headers,
		Body:     request.Body(),
		UniqueID: utils.Hash(request.Method() + request.URL().String() + request.Body()),
		Type:     RequestTypeNormal,
	}
}

func NewRequestFromDOM(path string, parent string) (*Request, error) {
	u, err := utils.ParseURL(path, parent)
	if err != nil {
		return nil, err
	}
	return &Request{
		URL:      u,
		Method:   MethodGET,
		Headers:  map[string]string{},
		Body:     "",
		UniqueID: utils.Hash(MethodGET + u.String() + ""),
		Type:     RequestTypeDOM,
	}, nil
}
