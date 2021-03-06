package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/sethgrid/curse"
)

const concurrency = 10
const displayRate = 200 * time.Millisecond

type progress struct {
	start         time.Time
	totalErrors   int
	totalRequests int
}

func main() {
	wg := sync.WaitGroup{}
	p := progress{
		start: time.Now(),
	}

	requestsChan := make(chan struct{})
	errorsChan := make(chan struct{})

	go func() {
		for {
			select {
			case <-requestsChan:
				p.totalRequests++
			case <-errorsChan:
				p.totalErrors++
			}
		}
	}()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			client := &http.Client{
				Transport: &http.Transport{
					MaxIdleConnsPerHost: 256,
				},
			}
			req, _ := http.NewRequest("GET", "http://echo.testing.eu.thermeon.io/", nil)
			for {
				resp, err := client.Do(req)
				requestsChan <- struct{}{}
				if err != nil {
					//fmt.Print("x")
					errorsChan <- struct{}{}
					continue
				}
				if body, err := ioutil.ReadAll(resp.Body); err != nil {
					errorsChan <- struct{}{}
					continue
				} else if len(body) != 21 {
					errorsChan <- struct{}{}
					continue
				}
				resp.Body.Close()

				// No real demand for this. Just a bit of foot off the gasishness to give some air in the pipeline
				time.Sleep(time.Millisecond)
			}
		}()
	}

	go display(&p)

	wg.Wait()
}

func display(p *progress) {

	c, err := curse.New()
	if err != nil {
		panic(err)
	}

	lastE, lastR := 0, 0

	for i := 0; i >= 0; i++ {
		diffR := p.totalRequests - lastR
		lastR = p.totalRequests
		diffE := p.totalErrors - lastE
		lastE = p.totalErrors

		totalSeconds := int(time.Since(p.start) / time.Second)
		elapsedSeconds := int(time.Second / displayRate)

		fmt.Printf("              Requests          Errors\n")
		fmt.Printf("Total:        %12d %12d\n", lastR, lastE)
		fmt.Printf("Current #/s:  %12d %12d\n", diffR*elapsedSeconds, diffE*elapsedSeconds)
		if totalSeconds > 0 {
			fmt.Printf("Average #/s:  %12d %12d\n", lastR/totalSeconds, lastE/totalSeconds)
		} else {
			fmt.Println("")
		}

		time.Sleep(displayRate)

		c.MoveUp(4)
	}

}
