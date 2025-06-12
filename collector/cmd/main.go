package main

import (
	"EGcollector/collector/collector"
	"sync"
)

func main() {
	director := collector.NewAutoFarmDirector()
	farmer := director.NewChromeFarmer()

	games := farmer.GetGames()
	wg := &sync.WaitGroup{}

	for _, game := range games {

		href, err := game.Attribute("href")
		if err != nil || "/en-US/free-games" == *href || href == nil {
			continue
		}

		wg.Add(1)
		go func(href string) {
			defer wg.Done()
			farmer.AddToCart(href)
		}(*href)

	}

	wg.Wait()
	farmer.Checkout()
}
