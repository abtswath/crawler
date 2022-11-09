package filter

import (
	"crawler/pkg/model"
	"strings"

	"github.com/go-rod/rod/lib/proto"
)

type DefaultFilter struct {
	host       string
	exclusions []string
}

func (d DefaultFilter) Can(request model.Request) bool {
	if request.URL.Host != d.host {
		return false
	}
	if d.shouldExclude(request) {
		return false
	}
	return request.ResourceType == proto.NetworkResourceTypeDocument
}

func (d DefaultFilter) shouldExclude(request model.Request) bool {
	for _, value := range d.exclusions {
		if strings.Contains(request.URL.Path, value) {
			return true
		}
	}
	return false
}

func New(host string, exclusions []string) Filter {
	return &DefaultFilter{
		host:       host,
		exclusions: exclusions,
	}
}
