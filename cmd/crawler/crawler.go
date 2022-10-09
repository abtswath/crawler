package main

import (
	"crawler/pkg/crawler"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"
)

func main() {
	u, _ := url.Parse("http://192.168.99.241/")
	c, err := crawler.New(*u, crawler.Option{
		Timeout: time.Second * 5,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()
	c.Run()
	result := c.Get()
	fmt.Println(result)
	fs, _ := os.OpenFile("./result.json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	defer fs.Close()
	json.NewEncoder(fs).Encode(result)
}
