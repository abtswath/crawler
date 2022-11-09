package filter

import "crawler/pkg/model"

type Filter interface {
	Can(model.Request) bool
}
