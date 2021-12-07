package main

import (
	"crawler"
	"fmt"
	"log"
	"net/url"
	"time"
)

func main() {
	target, _ := url.Parse("http://localhost:9003")
	c, err := crawler.New(crawler.Option{
		Timeout:     time.Second * 10,
		Incognito:   true,
		Headless:    false,
		Headers:     map[string]string{},
		PoolSize:    10,
		Target:      target,
		PageTimeout: time.Second * 5,
	})
	if err != nil {
		log.Fatalln(err)
	}
	c.Run()
	fmt.Println(c.Result)
	//
	//u := launcher.New().Headless(false).MustLaunch()
	//browser := rod.New().ControlURL(u).MustConnect()
	//defer browser.MustClose()
	//
	//// We create a pool that will hold at most 3 pages which means the max concurrency is 3
	//pool := rod.NewPagePool(10)
	//
	//// Create a page if needed
	//create := func() *rod.Page {
	//	// We use MustIncognito to isolate pages with each other
	//	return browser.MustPage()
	//}
	//
	//yourJob := func() {
	//	page := pool.Get(create)
	//	defer pool.Put(page)
	//
	//	page.MustNavigate("http://example.com").MustWaitLoad()
	//	fmt.Println(page.MustInfo().Title)
	//}
	//
	//// Run jobs concurrently
	//wg := sync.WaitGroup{}
	//for range "...................." {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		yourJob()
	//	}()
	//}
	//wg.Wait()
	//
	//// cleanup pool
	//pool.Cleanup(func(p *rod.Page) { p.MustClose() })
}
