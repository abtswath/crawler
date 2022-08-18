package task

import (
	"crawler/pkg/filter"
	"crawler/pkg/model"
	"encoding/json"
	"io"
)

type Collection struct {
	set    []model.Request
	filter filter.Filter
}

func (c *Collection) Put(req model.Request) {
	c.set = append(c.set, req)
}

func (c *Collection) Has(req model.Request) bool {
	return c.filter.Has(req)
}

func (c *Collection) GetAll() []model.Request {
	return c.set
}

func (c *Collection) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(c.set)
}
