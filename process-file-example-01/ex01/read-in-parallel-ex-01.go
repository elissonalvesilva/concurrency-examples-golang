package ex01

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/elissonalvesilva/golang-examples-concurrency/util"
)

type Model util.Model

func Run() {
	models := util.Open()

	const maxBatchSize int = 25
	skip := 0
	filesAmount := len(models)
	batchAmount := int(math.Ceil(float64(filesAmount / maxBatchSize)))
	var results []interface{}
	for i := 0; i <= batchAmount; i++ {
		lowerBound := skip
		upperBound := skip + maxBatchSize

		if upperBound > filesAmount {
			upperBound = filesAmount
		}

		batchItems := models[lowerBound:upperBound]

		skip += maxBatchSize

		processingErrorChan := make(chan error)
		processingDoneChan := make(chan int)
		processingErrors := make([]error, 0)

		go func() {
			for {
				select {
				case err := <-processingErrorChan:
					processingErrors = append(processingErrors, err)
				case <-processingDoneChan:
					close(processingErrorChan)
					close(processingDoneChan)
					return
				}
			}
		}()

		var itemProcessingGroup sync.WaitGroup
		itemProcessingGroup.Add(len(batchItems))

		for idx := range batchItems {
			go func(item *util.Model, idx int) {
				defer itemProcessingGroup.Done()
				results = append(results, item)
			}(&batchItems[idx], idx)
		}

		itemProcessingGroup.Wait()

		processingDoneChan <- 0

		if len(processingErrors) > 0 {
			fmt.Println("ERROR")
		}
	}
	fmt.Println(len(results))
	defer util.TimeTrack(time.Now(), "ex01")

}
