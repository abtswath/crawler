package utils

import (
	"crawler"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"net/url"
	"regexp"
	"strings"
)

func ShouldIgnoreRequest(request crawler.Request, keywords []string) bool {
	for _, keyword := range keywords {
		compile, err := regexp.Compile(keyword)
		if err != nil {
			continue
		}
		if compile.MatchString(request.URL.String()) {
			return true
		}
	}
	return false
}

func Hash(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func ParseURL(u string, parent ...string) (*url.URL, error) {
	return url.Parse(u)
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
		str.WriteString(randomStrSeed[rand.Intn(size):1])
	}
	return str.String()
}
