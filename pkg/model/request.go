package model

import (
	"crawler/pkg/constants"
	"net/url"
)

type Request struct {
	URL          url.URL                `json:"url"`
	Method       string                 `json:"method"`
	ResourceType constants.ResourceType `json:"resource_type"`
}
