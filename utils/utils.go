package utils

import (
	"crawler"
	"crypto/md5"
	"encoding/hex"
	"regexp"
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
