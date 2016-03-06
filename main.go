package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			j := 0
			for {
				if resp, err := http.Get("http://echo.testing.eu.thermeon.io/"); err != nil {
					fmt.Print("x")
				} else {
					if body, err := ioutil.ReadAll(resp.Body); err != nil {
						fmt.Print("e")
					} else if len(body) != 21 {
						fmt.Print("o")
					}
					resp.Body.Close()

					j = j + 1
					if j%100 == 0 {
						fmt.Print(".")
					}
				}
			}
		}()

	}
	wg.Wait()
}
