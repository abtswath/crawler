package collection

import (
	"crawler/request"
	"encoding/json"
	"github.com/go-rod/rod/lib/proto"
	"io"
	"sync"
)

type Item struct {
	URL          string                    `json:"url"`
	Method       string                    `json:"method"`
	Headers      map[string]string         `json:"headers"`
	Body         string                    `json:"body"`
	Type         int                       `json:"type"`
	ResourceType proto.NetworkResourceType `json:"resource_type"`
}

type Collection struct {
	list []*Item
	lock sync.Mutex
}

func (c *Collection) Put(req ...*request.Request) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, r := range req {
		c.list = append(c.list, &Item{
			URL:          r.URL.String(),
			Method:       r.Method,
			Headers:      r.Headers,
			Body:         r.Body,
			Type:         r.Type,
			ResourceType: r.ResourceType,
		})
	}
}

func (c *Collection) GetAll() []*Item {
	return c.list
}

func (c *Collection) JSON() ([]byte, error) {
	return json.Marshal(c.list)
}

func (c *Collection) Encode(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(c.list)
}
