package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, totalOfPizzas int

type Producer struct {
	pizzasOrderChan chan PizzaOrder
	quit            chan (chan error)
}

type PizzaOrder struct {
	id      int
	message string
	success bool
}

func (producer *Producer) Close() error {
	channel := make(chan error)
	producer.quit <- channel
	return <-channel

}

func makePizza(pizzaNumber int) *PizzaOrder {

	pizzaNumber++

	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order number #%d...\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		message := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}

		totalOfPizzas++

		fmt.Printf(
			"Making pizza number #%d. It will take %d seconds...\n",
			pizzaNumber, delay,
		)
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			message = fmt.Sprintf(
				"We ran out of ingredients for pizza #%d...", pizzaNumber,
			)
		} else if rnd <= 4 {
			message = fmt.Sprintf(
				"The cook had to leave while making pizza #%d...", pizzaNumber,
			)
		} else {
			success = true
			message = fmt.Sprintf(
				"Pizza #%d is ready!", pizzaNumber,
			)
		}

		pizzaOrder := PizzaOrder{
			id:      pizzaNumber,
			message: message,
			success: success,
		}

		return &pizzaOrder

	}

	pizzaOrder := PizzaOrder{
		id: pizzaNumber,
	}

	return &pizzaOrder

}

func pizzeriaRun(pizzaMaker *Producer) {
	var counter = 0
	for {
		currentPizza := makePizza(counter)

		counter = currentPizza.id

		select {
		//	Sending the pizza made to the pizza orders channel.
		case pizzaMaker.pizzasOrderChan <- *currentPizza:
			// Does nothing because we want to continue the loop
		//	Reading from quit channel to see if we should quit.
		case quitChannel := <-pizzaMaker.quit:
			close(pizzaMaker.pizzasOrderChan)
			close(quitChannel)
			return
		}
	}
}

func main() {

	rand.Seed(time.Now().UnixNano())

	color.Cyan("----------------------------")
	color.Cyan("Welcome to Vino's pizzeria")
	color.Cyan("----------------------------")

	pizzaJob := Producer{
		pizzasOrderChan: make(chan PizzaOrder),
		quit:            make(chan chan error),
	}

	go pizzeriaRun(&pizzaJob)

	//Consume pizza order
	for pizzaOrder := range pizzaJob.pizzasOrderChan {
		if pizzaOrder.id <= NumberOfPizzas {
			if pizzaOrder.success {
				color.Green("Order #%d was made successfully!", pizzaOrder.id)
			} else {
				color.Red("Order #%d failed!", pizzaOrder.id)
			}
		} else {
			color.Cyan("Done making pizzas for now...")

			err := pizzaJob.Close()

			if err != nil {
				panic(err.Error())
			}

		}
	}

	color.Cyan(
		"We made %d pizzas but failed to make %d in a total of %d...",
		pizzasMade, pizzasFailed, NumberOfPizzas,
	)

}
