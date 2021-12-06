package filter

import (
	"crawler/request"
	"github.com/emirpasic/gods/sets/hashset"
)

type Filter interface {
	Exists(request *request.Request) bool

	Static(request *request.Request) bool
}

type DefaultFilter struct {
	set *hashset.Set
}

func NewDefaultFilter() *DefaultFilter {
	return &DefaultFilter{
		set: hashset.New(),
	}
}

func (d *DefaultFilter) Exists(request *request.Request) bool {
	if !d.set.Contains(request.UniqueID) {
		d.set.Add(request.UniqueID)
		return false
	}
	return true
}

func (d *DefaultFilter) Static(request *request.Request) bool {
	return false
}
