package main

import (
	"EGcollector/collector/collector"
	"sync"
)

func main() {
	director := collector.NewCollectorDirector()
	collector := director.NewChromeCollector()

	games := collector.GetGames()
	wg := &sync.WaitGroup{}

	for _, game := range games {

		href, err := game.Attribute("href")
		if err != nil || "/en-US/free-games" == *href || href == nil {
			continue
		}

		wg.Add(1)
		go func(href string) {
			defer wg.Done()
			collector.AddToCart(href)
		}(*href)

	}

	wg.Wait()
	collector.Checkout()
}
