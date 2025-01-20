package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Doehnert/currency/internal/currency"
)

func runCurrencyWorker(
	workerId int,
	currencyChan <-chan currency.Currency,
	resultChan chan<- currency.Currency) {
	fmt.Printf("Worker %d started\n", workerId)
	for c := range currencyChan {
		err := c.FetchCurrencyRate()
		if err != nil {
			panic(err)
		}
		resultChan <- c
	}
	fmt.Printf("Worker %d finished\n", workerId)
}

func main() {

	ce := &currency.MyCurrencyExchange{
		Currencies: make(map[string]currency.Currency),
	}

	err := ce.FetchAllCurrencies()
	if err != nil {
		panic(err)
	}

	currencyChan := make(chan currency.Currency, len(ce.Currencies))
	resultChan := make(chan currency.Currency, len(ce.Currencies))

	numOfWorkers := runtime.NumCPU()
	fmt.Println("Number of workers: ", numOfWorkers)

	for i := 0; i < numOfWorkers*100; i++ {
		go runCurrencyWorker(i, currencyChan, resultChan)
	}

	startTime := time.Now()

	resultCount := 0

	for _, c := range ce.Currencies {
		currencyChan <- c
	}

	for {
		if resultCount == len(ce.Currencies) {
			fmt.Println("Closing resultChan")
			close(currencyChan)
			break
		}
		select {
		case c := <-resultChan:
			ce.Currencies[c.Code] = c
			resultCount++
		case <-time.After(3 * time.Second):
			fmt.Println("Timeout")
			return
		}
	}

	endTime := time.Now()

	fmt.Println("---------- Results ----------")
	for code, c := range ce.Currencies {
		fmt.Printf("%s (%s): %d ranges\n", c.Name, code, len(c.Rates))
		fmt.Println("----------------------------")
		fmt.Println("Time taken: ", endTime.Sub(startTime))
	}
}
