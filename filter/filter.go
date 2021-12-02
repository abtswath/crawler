package filter

import "crawler"

type Filter interface {
	Exists(request *crawler.Request) bool

	Static(request *crawler.Request) bool
}

type DefaultFilter struct {
}

func (d *DefaultFilter) Exists(request *crawler.Request) bool {
	return false
}

func (d *DefaultFilter) Static(request *crawler.Request) bool {
	return false
}
