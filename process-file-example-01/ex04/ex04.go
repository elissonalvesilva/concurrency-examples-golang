package ex04

import (
	"fmt"
	"sync"
	"time"

	"github.com/elissonalvesilva/golang-examples-concurrency/util"
)

type Model util.Model

/*
 * Fan out, fan in pattern
 */

func Run() {
	models := util.Open()
	var resultsFinal []util.Model

	in := gen(models)

	// Distribute the sq work across two goroutines that both read from in.
	model1 := sq(in)
	model2 := sq(in)

	// Consume the merged output from model1 and model2.
	for m := range merge(model1, model2) {
		resultsFinal = append(resultsFinal, m)
	}

	fmt.Println(len(resultsFinal))
	defer util.TimeTrack(time.Now(), "ex04-Fan out, fan in pattern")

}

func gen(models []util.Model) <-chan util.Model {
	out := make(chan util.Model)
	go func() {
		for _, v := range models {
			out <- v
		}
		close(out)
	}()

	return out
}

func sq(model <-chan util.Model) <-chan util.Model {
	out := make(chan util.Model)

	go func() {
		for v := range model {
			if v.Type == "CreateEvent" {
				out <- v
			}
		}
		close(out)
	}()

	return out
}

func merge(models ...<-chan util.Model) <-chan util.Model {
	var wg sync.WaitGroup

	out := make(chan util.Model)
	// Start an output goroutine for each input channel in models.  output
	// copies values from c to out until m is closed, then calls wg.Done.
	output := func(model <-chan util.Model) {
		for m := range model {
			out <- m
		}
		wg.Done()
	}

	wg.Add(len(models))
	for _, m := range models {
		go output(m)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
