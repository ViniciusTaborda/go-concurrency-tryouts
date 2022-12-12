package main

import (
	"fmt"
	"sync"
)

var msg string

func UpdateMessage(text string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	msg = text
}

func PrintMessage() {
	fmt.Println(msg)
}

func main() {

	// challenge: modify this code so that the calls to UpdateMessage() on lines
	// 28, 30, and 33 run as goroutines, and implement wait groups so that
	// the program runs properly, and prints out three different messages.
	// Then, write a test for all three functions in this program: UpdateMessage(),
	// PrintMessage(), and main().

	msg = "Hello, world!"

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go UpdateMessage("Hello, universe!", &waitGroup)
	waitGroup.Wait()
	PrintMessage()

	waitGroup.Add(1)
	go UpdateMessage("Hello, cosmos!", &waitGroup)
	waitGroup.Wait()
	PrintMessage()

	waitGroup.Add(1)
	go UpdateMessage("Hello, world!", &waitGroup)
	waitGroup.Wait()
	PrintMessage()

}
