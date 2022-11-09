package utils

import (
	"testing"
)

func TestParseURL(t *testing.T) {
	parentURL := "https://baidu.com/qwe/asd?a=1#zxc"
	u, err := ParseURL("/abc", parentURL)
	if err != nil {
		t.Fatalf("parse error: %v\n", err)
		t.Fail()
	}
	t.Logf("/abc url: %v\n", u)
	if u.String() != "https://baidu.com/abc" {
		t.Fatalf("/abc url string: %s\n", u.String())
		t.Fail()
	}

	u, err = ParseURL("abc", parentURL)
	if err != nil {
		t.Fatalf("parse error: %v\n", err)
		t.Fail()
	}
	t.Logf("abc url: %v\n", u)
	if u.String() != "https://baidu.com/qwe/asd/abc" {
		t.Fatalf("abc url string: %s\n", u.String())
		t.Fail()
	}

	u, err = ParseURL("https://baidu.com/abc", parentURL)
	if err != nil {
		t.Fatalf("parse error: %v\n", err)
		t.Fail()
	}
	if u.String() != "https://baidu.com/abc" {
		t.Fatalf("https://baidu.com/abc url string: %s\n", u.String())
		t.Fail()
	}
}
