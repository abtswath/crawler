package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"path"
	"strings"
)

func Hash(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func ParseURL(u string, parent string) (*url.URL, error) {
	u = strings.Trim(u, " ")

	if len(u) == 0 {
		return nil, errors.New("invalid url")
	}
	if len(parent) <= 0 {
		return url.Parse(u)
	}
	if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		return url.Parse(u)
	}
	if strings.HasPrefix(u, "javascript:") || strings.HasPrefix(u, "mailto:") {
		return nil, errors.New("invalid url")
	}
	parentURL, err := url.Parse(parent)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(u, "/") {
		u = path.Join(parentURL.Path, u)
	}
	return url.Parse(fmt.Sprintf("%s://%s%s", parentURL.Scheme, parentURL.Host, u))
}

func StringArrayInclude(haystack []string, value string) bool {
	for _, s := range haystack {
		if s == value {
			return true
		}
	}
	return false
}

var randomStrSeed = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomStr(length int) string {
	size := len(randomStrSeed)
	str := strings.Builder{}
	for i := 0; i < length; i++ {
		index := rand.Intn(size)
		str.WriteString(randomStrSeed[index : index+1])
	}
	return str.String()
}
