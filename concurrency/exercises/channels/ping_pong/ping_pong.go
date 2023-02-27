package main

import (
	"fmt"
	"strings"
)

// Here declaring ping as a send only chan and pong as a receive only chan
func shout(ping <-chan string, pong chan<- string) {
	for {
		message := <-ping
		pong <- fmt.Sprintf("%s !!!", strings.ToUpper(message))
	}
}

func main() {

	ping := make(chan string)
	pong := make(chan string)

	go shout(ping, pong)

	fmt.Println("Type something and press ENTER. ")
	fmt.Println("(Empty input to quit) ")
	fmt.Println()

	for {
		fmt.Print("-> ")

		var userInput string

		_, _ = fmt.Scanln(&userInput)

		if strings.ToLower(userInput) == "" {
			break
		}

		ping <- userInput

		response := <-pong

		fmt.Printf("Response is %s.\n", response)
	}

	fmt.Println("All done. Closing all channels...")

	close(ping)
	close(pong)

}
