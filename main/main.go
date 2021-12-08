package main

import (
	"crawler"
	"fmt"
	"log"
	"net/url"
	"time"
)

func main() {
	target, _ := url.Parse("http://localhost:94/v2")
	c, err := crawler.New(crawler.Option{
		Timeout:   time.Minute * 20,
		Incognito: true,
		Headless:  true,
		Headers: map[string]string{
			"From": "Crawler",
		},
		PoolSize:       10,
		Target:         target,
		PageTimeout:    time.Second * 5,
		BrowserTrace:   false,
		IgnoreKeywords: []string{"delete", "remove", "Remove", "Delete", "logout", "exit"},
	})
	if err != nil {
		log.Fatalln(err)
	}
	c.Run()
	for _, request := range c.Result {
		fmt.Printf("Method: %s, URL: %s\n", request.Method, request.URL.String())
	}
}
