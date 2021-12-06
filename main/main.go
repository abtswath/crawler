package main

import (
	"crawler"
	"crawler/config"
	"fmt"
	"log"
	"net/url"
)

func main() {
	target, _ := url.Parse("http://localhost:9003")
	c, err := crawler.NewCrawler(config.NewOption(target))
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()
	c.Run()
	fmt.Println(c.Result)
}
