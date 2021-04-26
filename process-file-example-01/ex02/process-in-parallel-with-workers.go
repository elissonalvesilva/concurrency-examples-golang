package ex02

import (
	"fmt"
	"sync"
	"time"

	"github.com/elissonalvesilva/golang-examples-concurrency/util"
)

type Model util.Model

func Run() {
	models := util.Open()
	var results = make(chan util.Model)
	var wg sync.WaitGroup
	var resultsFinal []util.Model

	for _, v := range models {
		wg.Add(1)
		go worker(v, &wg, results)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	for v := range results {
		resultsFinal = append(resultsFinal, v)
	}

	fmt.Println(len(models))
	fmt.Println(len(resultsFinal))
	defer util.TimeTrack(time.Now(), "ex02")
}

// ReleaseEvent
// PushEvent
// IssuesEvent
// WatchEvent
// CreateEvent
// PullRequestEvent
func worker(model util.Model, wg *sync.WaitGroup, results chan util.Model) {
	defer wg.Done()

	results <- model

}
