package main

import (
	"fmt"
	"time"
)

func listenToChan(ch chan int) {
	for {
		i := <-ch
		fmt.Printf("Got %d from channel!\n", i)

		//Simulate some logic
		time.Sleep(1 * time.Second)
	}
}

func main() {

	//Buffering channel...
	ch := make(chan int, 10)

	go listenToChan(ch)

	for i := 0; i < 100; i++ {
		fmt.Printf("Sending %d to channel!\n", i)
		ch <- i
		fmt.Printf("Sent %d to channel!\n", i)
	}

	fmt.Println("All done!")
	close(ch)
}
