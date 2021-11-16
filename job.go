package crawler

func (c *Crawler) job(target string) {
	defer c.wg.Done()
	page := c.pool.Get(c.newPage)
	defer c.pool.Put(page)
	err := page.Navigate(target)
	if err != nil {
		return
	}
	err = page.WaitLoad()
	if err != nil {
		return
	}
	// TODO. Collect URL
	// TODO. Fire events
	// TODO. Fill forms
}

func (c *Crawler) newJob(target string) {
	c.wg.Add(1)
	go c.job(target)
}
