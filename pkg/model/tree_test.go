package model

import (
	"net/url"
	"testing"

	"github.com/go-rod/rod/lib/proto"
)

func TestTreeAddAndGet(t *testing.T) {
	node := &Node{}

	paths := []string{
		"/",
		"/cmd/test/",
		"/cmd/test/3",
		"/cmd/who/",
		"/cmd/whoami",
		"/cmd/whoami/r",
		"/cmd/whoami/root",
		"/cmd/whoami/root/",
		"/src/",
		"/src/some/file.png",
		"/search/",
		"/search/someth!ng+in+ünìcodé",
		"/search/gin",
		"/search/gin-gonic",
		"/search/google",
		"/user_gopher",
		"/user_gopher/about",
		"/files/js/inc/framework.js",
		"/info/gordon/public",
		"/info/gordon/project/go",
		"/info/gordon/project/golang",
		"/aa/aa",
		"/ab/ab",
		"/a",
		"/all",
		"/d",
		"/ad",
		"/dd",
		"/dddaa",
		"/aa",
		"/aaa",
		"/aaa/cc",
		"/ab",
		"/abb",
		"/abb/cc",
		"/allxxxx",
		"/alldd",
		"/all/cc",
		"/a/cc",
		"/c1/d/e",
		"/c1/d/e1",
		"/c1/d/ee",
		"/cc/cc",
		"/ccc/cc",
		"/deedwjfs/cc",
		"/acllcc/cc",
		"/get/test/abc/",
		"/get/te/abc/",
		"/get/testaa/abc/",
		"/get/xx/abc/",
		"/get/tt/abc/",
		"/get/a/abc/",
		"/get/t/abc/",
		"/get/aa/abc/",
		"/get/abas/abc/",
		"/something/secondthing/test",
		"/something/abcdad/thirdthing",
		"/something/secondthingaaaa/thirdthing",
		"/something/se/thirdthing",
		"/something/s/thirdthing",
		"/c/d/ee",
		"/c/d/e/ff",
		"/c/d/e/f/gg",
		"/c/d/e/f/g/hh",
		"/cc/dd/ee/ff/gg/hh",
		"/get/abc",
		"/get/a",
		"/get/abz",
		"/get/12a",
		"/get/abcd",
		"/get/abc/123abc",
		"/get/abc/12",
		"/get/abc/123ab",
		"/get/abc/xyz",
		"/get/abc/123abcddxx",
		"/get/abc/123abc/xxx8",
		"/get/abc/123abc/x",
		"/get/abc/123abc/xxx",
		"/get/abc/123abc/abc",
		"/get/abc/123abc/xxx8xxas",
		"/get/abc/123abc/xxx8/1234",
		"/get/abc/123abc/xxx8/1",
		"/get/abc/123abc/xxx8/123",
		"/get/abc/123abc/xxx8/78k",
		"/get/abc/123abc/xxx8/1234xxxd",
		"/get/abc/123abc/xxx8/1234/ffas",
		"/get/abc/123abc/xxx8/1234/f",
		"/get/abc/123abc/xxx8/1234/ffa",
		"/get/abc/123abc/xxx8/1234/kka",
		"/get/abc/123abc/xxx8/1234/ffas321",
		"/get/abc/123abc/xxx8/1234/kkdd/12c",
		"/get/abc/123abc/xxx8/1234/kkdd/1",
		"/get/abc/123abc/xxx8/1234/kkdd/12",
		"/get/abc/123abc/xxx8/1234/kkdd/12b",
		"/get/abc/123abc/xxx8/1234/kkdd/34",
		"/get/abc/123abc/xxx8/1234/kkdd/12c2e3",
		"/get/abc/12/test",
		"/get/abc/123abdd/test",
		"/get/abc/123abdddf/test",
		"/get/abc/123ab/test",
		"/get/abc/123abgg/test",
		"/get/abc/123abff/test",
		"/get/abc/123abffff/test",
		"/get/abc/123abd/test",
		"/get/abc/123abddd/test",
		"/get/abc/123/test22",
		"/get/abc/123abg/test",
		"/get/abc/123abf/testss",
		"/get/abc/123abfff/te",
		"/hi",
		"/contact",
		"/co",
		"/con",
		"/cona",
		"/no",
		"/ab",
		"/α",
		"/β",
	}

	for _, path := range paths {
		node.Put(path, url.URL{Path: path}, proto.NetworkResourceTypeDocument)
	}

	for _, path := range paths {
		if node.Get(path) != nil {
			t.Errorf("missing path '%s'\n", path)
			t.FailNow()
		} else {
			n := node.Get(path)
			if n.URL.Path != path {
				t.Errorf("mismatch for path '%s': got '%s'\n", path, n.URL.Path)
			}
		}
	}
}
