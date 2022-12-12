package main

import (
	"fmt"
	"sync"
)

var waitGroup sync.WaitGroup

type Income struct {
	Source string
	Amount float64
}

func main() {

	var bankBalance float64
	var mutex sync.Mutex

	incomes := []Income{
		{Source: "Main job", Amount: 500.23},
		{Source: "Gifts", Amount: 50.99},
		{Source: "Part time job", Amount: 15.50},
		{Source: "Investments", Amount: 100},
	}

	waitGroup.Add(len(incomes))

	for index, income := range incomes {
		go func(index int, income Income) {

			defer waitGroup.Done()

			for week := 1; week <= 52; week++ {
				mutex.Lock()

				temp := bankBalance
				temp += income.Amount
				bankBalance = temp

				mutex.Unlock()

				fmt.Printf(
					"On week %d, you earned $%f from %s! \n",
					week,
					income.Amount,
					income.Source,
				)
			}

		}(index, income)
	}

	waitGroup.Wait()

	fmt.Printf("Final bank balance: %f", bankBalance)

}
