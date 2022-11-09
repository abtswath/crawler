package model

import (
	"net/url"

	"github.com/go-rod/rod/lib/proto"
)

type Request struct {
	URL          url.URL
	Method       string
	ResourceType proto.NetworkResourceType
}

type Result struct {
	URL          string                    `json:"url"`
	Method       string                    `json:"method"`
	ResourceType proto.NetworkResourceType `json:"resource_type"`
}
