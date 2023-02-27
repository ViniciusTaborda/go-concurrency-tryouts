package main

import (
	"fmt"
	"sync"
)

func printStringln(text string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	fmt.Println(text)
}

func main() {

	// Wait groups are the easiest way to sync multiple goroutines
	var waitGroup sync.WaitGroup

	words := []string{
		"alpha",
		"beta",
		"delta",
		"gamma",
		"pi",
		"zeta",
		"eta",
		"theta",
		"episolon",
	}

	waitGroup.Add(len(words))

	for index, word := range words {
		go printStringln(fmt.Sprintf("%d: %s", index, word), &waitGroup)
	}

	waitGroup.Wait()
}
