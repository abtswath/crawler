package filter

import (
	"crawler/request"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-rod/rod/lib/proto"
)

type Filter interface {
	Allow(r *request.Request) bool

	Exists(r *request.Request) bool

	Static(r *request.Request) bool
}

type DefaultFilter struct {
	set      *hashset.Set
	RootHost string
}

func NewDefaultFilter() *DefaultFilter {
	return &DefaultFilter{
		set: hashset.New(),
	}
}

func (d *DefaultFilter) Allow(r *request.Request) bool {
	return r.URL.Host == d.RootHost
}

func (d *DefaultFilter) Exists(r *request.Request) bool {
	if !d.set.Contains(r.UniqueID) {
		d.set.Add(r.UniqueID)
		return false
	}
	return true
}

func (d *DefaultFilter) Static(r *request.Request) bool {
	return r.ResourceType != proto.NetworkResourceTypeDocument
}
