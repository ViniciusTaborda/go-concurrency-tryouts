package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Barbershop struct {
	NumberOfSeats     int
	HaircutDuration   time.Duration
	NumberOfBarbers   int
	BarberDoneChannel (chan bool)
	ClientsChannel    (chan string)
	IsOpen            bool
}

func (barbershop *Barbershop) Close() {

	color.Blue("Closing shop for the day!")

	close(barbershop.ClientsChannel)
	barbershop.IsOpen = false

	for barber := 1; barber <= barbershop.NumberOfBarbers; barber++ {
		<-barbershop.BarberDoneChannel
	}

	close(barbershop.BarberDoneChannel)

	fmt.Println()
	color.Blue("Barber shop is closed and everyone went home...")

}

func (barbershop *Barbershop) AddClient(clientName string) {

	color.Cyan(" - %s arrived!", clientName)

	if barbershop.IsOpen {
		select {
		case barbershop.ClientsChannel <- clientName:
			color.Cyan(" - %s takes a set to wait for a cut...", clientName)
		default:
			color.Red(" - The waiting room is full, so %s leaves...", clientName)
		}
	} else {
		color.Red(" - Barbershop is closed, so %s leaves...", clientName)
	}

}

func (barbershop *Barbershop) CutHair(barberName, clientName string) {

	color.Green("%s is cutting %s's hair!", barberName, clientName)
	time.Sleep(barbershop.HaircutDuration)
	color.Green("%s finished cutting %s's hair!", barberName, clientName)

}

func (barbershop *Barbershop) SendBarberHome(barberName string) {
	color.Red("%s is done for the day...", barberName)

	// Sending a signal that this barber is done for the day
	barbershop.BarberDoneChannel <- true
}

func (barbershop *Barbershop) AddBarber(barberName string) {
	barbershop.NumberOfBarbers++

	go func() {
		isSleeping := false

		color.Cyan(
			"%s goes to the waiting room to check for clients...",
			barberName,
		)

		for {
			// If there are no clients waiting to get a cut, takes a nap
			if len(barbershop.ClientsChannel) == 0 {
				color.Cyan(
					"%s sees that there are not clients, takes a nap "+
						"ZZzzzZZzzzZZzzzz...",
					barberName,
				)
				isSleeping = true
			}
			// The second return of reading a channel indicates
			// if the channel is open or not.
			clientName, barbershopOpen := <-barbershop.ClientsChannel

			if barbershopOpen {
				if isSleeping {
					color.Cyan(
						"%s wakes %s up!",
						clientName,
						barberName,
					)
					isSleeping = false
				}

				// Not sleeping so cuts hair.
				barbershop.CutHair(barberName, clientName)

			} else {
				// Shop is closed, go home.
				barbershop.SendBarberHome(barberName)

				// Close go routine, return function
				return
			}

		}

	}()

}
