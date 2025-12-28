package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var images = map[string]string{
	"imgIphone":    "https://images.unsplash.com/photo-1695048133142-1a20484d2569?q=80&w=800&auto=format&fit=crop",
	"imgHeadphone": "https://images.unsplash.com/photo-1618366712010-f4ae9c647dcb?q=80&w=800&auto=format&fit=crop",
	"imgTshirt":    "https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?q=80&w=800&auto=format&fit=crop",
	"imgJacket":    "https://images.unsplash.com/photo-1523275335684-37898b6baf30?q=80&w=800&auto=format&fit=crop",
	"imgSalad":     "https://images.unsplash.com/photo-1546069901-ba9599a7e63c?q=80&w=800&auto=format&fit=crop",
	"imgFruit":     "https://images.unsplash.com/photo-1610832958506-aa56368176cf?q=80&w=800&auto=format&fit=crop",
	"imgSeafood":   "https://images.unsplash.com/photo-1534483509719-3feaee7c30da?q=80&w=800&auto=format&fit=crop",
	"imgFridge":    "https://images.unsplash.com/photo-1588854337221-4cf9fa96059c?q=80&w=800&auto=format&fit=crop",
	"imgWasher":    "https://images.unsplash.com/photo-1626806819282-2c1dc01a5e0c?q=80&w=800&auto=format&fit=crop",
	"imgLamp":      "https://images.unsplash.com/photo-1513506003011-3b03c860ad80?q=80&w=800&auto=format&fit=crop",
	"imgKeyboard":  "https://images.unsplash.com/photo-1595225476474-87563907a212?q=80&w=800&auto=format&fit=crop",
	"imgShoes":     "https://images.unsplash.com/photo-1542291026-7eec264c27ff?q=80&w=800&auto=format&fit=crop",
}

func main() {
	var wg sync.WaitGroup
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for name, url := range images {
		wg.Add(1)
		go func(n, u string) {
			defer wg.Done()
			resp, err := client.Head(u)
			if err != nil {
				fmt.Printf("❌ %s: Error %v\n", n, err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				fmt.Printf("❌ %s: Status %d\n", n, resp.StatusCode)
			} else {
				fmt.Printf("✅ %s: OK\n", n)
			}
		}(name, url)
	}

	wg.Wait()
}
