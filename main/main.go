package main

import (
	"crawler"
	"github.com/go-rod/rod/lib/proto"
	"github.com/sirupsen/logrus"
	"log"
	"net/url"
	"os"
	"time"
)

func main() {
	target, _ := url.Parse("http://localhost:94/v2")
	c, err := crawler.New(crawler.Option{
		Timeout:   time.Minute,
		Incognito: true,
		Headless:  false,
		Headers: map[string]string{
			"From": "Crawler",
		},
		PoolSize:       20,
		Target:         target,
		PageTimeout:    time.Second * 5,
		BrowserTrace:   false,
		IgnoreKeywords: []string{"delete", "remove", "Remove", "Delete", "logout", "exit"},
		UploadFile:     "./image.png",
		Cookies: []*proto.NetworkCookieParam{
			{
				Name:   "PHPSESSID",
				Value:  "9pi8ttm1ome4qe5g31fp7veid2",
				Domain: "localhost",
				Path:   "/",
			},
			{
				Name:   "defaf8b50191c3f357502de1967c4266_apiconfig",
				Value:  "1",
				Domain: "localhost",
				Path:   "/",
			},
			{
				Name:   "defaf8b50191c3f357502de1967c4266_enhance",
				Value:  "1",
				Domain: "localhost",
				Path:   "/",
			},
			{
				Name:   "defaf8b50191c3f357502de1967c4266_firewallsetting",
				Value:  "1",
				Domain: "localhost",
				Path:   "/",
			},
			{
				Name:   "defaf8b50191c3f357502de1967c4266_ip_range",
				Value:  "*.*.*.*",
				Domain: "localhost",
				Path:   "/",
			},
			{
				Name:   "defaf8b50191c3f357502de1967c4266_ip_range_user",
				Value:  "*.*.*.*",
				Domain: "localhost",
				Path:   "/",
			},
			{
				Name:   "defaf8b50191c3f357502de1967c4266_linkageconfig",
				Value:  "1",
				Domain: "localhost",
				Path:   "/",
			},
			{
				Name:   "defaf8b50191c3f357502de1967c4266_loginName",
				Value:  "webadmin",
				Domain: "localhost",
				Path:   "/",
			},
			{
				Name:   "defaf8b50191c3f357502de1967c4266_max_host_thread_global",
				Value:  "20",
				Domain: "localhost",
				Path:   "/",
			},
			{
				Name:   "defaf8b50191c3f357502de1967c4266_max_port_thread_global",
				Value:  "10",
				Domain: "localhost",
				Path:   "/",
			},
			{
				Name:   "defaf8b50191c3f357502de1967c4266_max_weak_thread_global",
				Value:  "20",
				Domain: "localhost",
				Path:   "/",
			},
			{
				Name:   "defaf8b50191c3f357502de1967c4266_max_web_thread_global",
				Value:  "10",
				Domain: "localhost",
				Path:   "/",
			},
			{
				Name:   "defaf8b50191c3f357502de1967c4266_shouldModifyPassword",
				Value:  "0",
				Domain: "localhost",
				Path:   "/",
			},
			{
				Name:   "defaf8b50191c3f357502de1967c4266_type",
				Value:  "2",
				Domain: "localhost",
				Path:   "/",
			},
			{
				Name:   "defaf8b50191c3f357502de1967c4266_upgradesetting",
				Value:  "1",
				Domain: "localhost",
				Path:   "/",
			},
		},
		LogLevel: logrus.TraceLevel,
	})
	if err != nil {
		log.Fatalln(err)
	}
	c.Run()
	resultFile, err := os.OpenFile("./result.json", os.O_CREATE|os.O_TRUNC|os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	defer resultFile.Close()
	err = c.Result.Encode(resultFile)
	if err != nil {
		log.Fatalln(err)
	}
}
