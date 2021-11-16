package crawler

import "net/url"

type URL struct {
	*url.URL
	Method      string
}
