package main

import (
	"EGcollector/collector/collector"
	"sync"
)

func main() {
	director := collector.NewCollectorDirector()
	collector := director.NewChromeCollector()
	wg := sync.WaitGroup{}

	games := collector.GetGames()

	for _, game := range games {
		_, err := game.Element(".css-gyjcm9")
		if err == nil {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			href, err := game.Attribute("href")
			if err != nil || "/en-US/free-games" == *href || href == nil {
				return
			}

			collector.AddToCart(*href)
		}()

	}
	wg.Wait()
	collector.Checkout()
}
