package filter

import "crawler/pkg/model"

type Filter interface {
	Has(request model.Request) bool
	Can(request model.Request) bool
}
