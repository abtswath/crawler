package utils

import (
	"errors"
	"net/url"
	"strings"
)

func ParseURL(str string, urls ...string) (*url.URL, error) {
	u, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	if len(urls) <= 0 && u.Host == "" {
		return nil, errors.New("invalid url")
	}
	parentURL, err := url.Parse(urls[0])
	if err != nil {
		return nil, err
	}
	u.Host = parentURL.Host
	u.User = parentURL.User
	u.Scheme = parentURL.Scheme
	if !strings.HasPrefix(u.Path, "/") {
		prefix := parentURL.Path
		if strings.HasSuffix(parentURL.Path, "/") {
			prefix = strings.TrimRight(prefix, "/")
		}
		u.Path = strings.Join([]string{prefix, str}, "/")
	}
	return u, nil
}

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func Min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}
