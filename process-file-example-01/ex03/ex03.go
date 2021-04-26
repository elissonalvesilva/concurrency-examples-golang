package ex03

import (
	"fmt"
	"time"

	"github.com/elissonalvesilva/golang-examples-concurrency/util"
)

type Model util.Model

/*
 * Squaring Numbers Patterns
 */

func Run() {
	models := util.Open()
	var resultsFinal []util.Model

	c := gen(models)
	out := sq(c)

	for v := range out {
		resultsFinal = append(resultsFinal, v)
	}
	fmt.Println(len(resultsFinal))
	defer util.TimeTrack(time.Now(), "ex03-squaring numbers patterns")

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
