package main

import (
	"crawler/pkg/crawler"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/go-rod/rod/lib/proto"
)

func main() {
	u, _ := url.Parse("https://192.168.99.114/v2")
	c, err := crawler.New(*u, crawler.Option{
		Timeout: time.Second * 5,
		Cookies: []*proto.NetworkCookieParam{
			{
				Name:   "ce5c5c357ab0c54a7f75e0f8d13221c1_loginName",
				Value:  "webadmin",
				Domain: "192.168.99.114",
				Path:   "/",
			},
			{
				Name:   "nvssession",
				Value:  "0fnasfmeqe6r7q369c070hp4b1",
				Domain: "192.168.99.114",
				Path:   "/",
			},
			{
				Name:   "ce5c5c357ab0c54a7f75e0f8d13221c1_type",
				Value:  "2",
				Domain: "192.168.99.114",
				Path:   "/",
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()
	c.Run()
	result := c.Get()
	fmt.Println(result)
}
