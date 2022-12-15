package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// This is a simple demonstration of how to solve the Sleeping Barber dilemma, a classic computer science problem
// which illustrates the complexities that arise when there are multiple operating system processes. Here, we have
// a finite number of barbers, a finite number of seats in a waiting room, a fixed length of time the barbershop is
// open, and clients arriving at (roughly) regular intervals. When a barber has nothing to do, he or she checks the
// waiting room for new clients, and if one or more is there, a haircut takes place. Otherwise, the barber goes to
// sleep until a new client arrives. So the rules are as follows:
//
//		- if there are no customers, the barber falls asleep in the chair
//		- a customer must wake the barber if he is asleep
//		- if a customer arrives while the barber is working, the customer leaves if all chairs are occupied and
//		  sits in an empty chair if it's available
//		- when the barber finishes a haircut, he inspects the waiting room to see if there are any waiting customers
//		  and falls asleep if there are none
// 		- shop can stop accepting new clients at closing time, but the barbers cannot leave until the waiting room is
//	      empty
//		- after the shop is closed and there are no clients left in the waiting area, the barber
//		  goes home
//
// The Sleeping Barber was originally proposed in 1965 by computer science pioneer Edsger Dijkstra.
//
// The point of this problem, and its solution, was to make it clear that in a lot of cases, the use of
// semaphores (mutexes) is not needed.

var numberOfSeats = 2
var arrivalRate = 100
var haircutDuration = 1000 * time.Millisecond
var barbershopOpenDuration = 10 * time.Second

func main() {

	// Seed our random number generator
	rand.Seed(time.Now().UnixNano())

	// Print welcome message
	color.Yellow(" -=-=-=- The sleeping barber problem -=-=-=- ")
	fmt.Println()

	// Create channels
	clientsChannel := make(chan string, numberOfSeats)
	barberDoneChannel := make(chan bool)

	// Create some data structure that represents the barbershop

	barbershop := Barbershop{
		NumberOfSeats:     numberOfSeats,
		HaircutDuration:   haircutDuration,
		NumberOfBarbers:   0,
		BarberDoneChannel: barberDoneChannel,
		ClientsChannel:    clientsChannel,
		IsOpen:            true,
	}

	color.Green("The shop is open for the day...")
	fmt.Println()

	// Add barbers (consumers)
	barbershop.AddBarber("Jovi")
	barbershop.AddBarber("Vino")

	// Run the barbershop as a goroutine
	shopClosingChannel := make(chan bool)
	closedShopChannel := make(chan bool)

	go func() {
		<-time.After(barbershopOpenDuration)
		shopClosingChannel <- true
		barbershop.Close()
		closedShopChannel <- true
	}()

	// Add clients (producers)

	clientId := 0

	go func() {
		for {
			// Get random number with average arrival rate
			randomMiliSeconds := rand.Int() % (2 * arrivalRate)

			select {
			case <-shopClosingChannel:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMiliSeconds)):
				barbershop.AddClient(fmt.Sprintf("Client #%d", clientId))
				clientId++
			}

		}
	}()

	// Wait until the barbershop is closed
	<-closedShopChannel
}
