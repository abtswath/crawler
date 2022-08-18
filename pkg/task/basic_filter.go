package task

import (
	"crawler/pkg/model"
)

type BasicFilter struct {
	domains    []string
	collection *Collection
}

func (b *BasicFilter) Has(request model.Request) bool {
	urlStr := request.URL.String()
	for _, req := range b.collection.GetAll() {
		if req.URL.String() == urlStr && req.Method == request.Method {
			return true
		}
	}
	return false
}

func (b *BasicFilter) Can(request model.Request) bool {
	if b.domains == nil || len(b.domains) <= 0 {
		return true
	}
	host := request.URL.Hostname()
	for _, domain := range b.domains {
		if domain == host {
			return true
		}
	}
	return false
}
