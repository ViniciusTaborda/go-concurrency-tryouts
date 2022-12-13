package main

import (
	"fmt"
	"sync"
	"time"
)

// The Dining Philosophers problem is well known in computer science circles.
// Five philosophers, numbered from 0 through 4, live in a house where the
// table is laid for them; each philosopher has their own place at the table.
// Their only difficulty – besides those of philosophy – is that the dish
// served is a very difficult kind of spaghetti which has to be eaten with
// two forks. There are two forks next to each plate, so that presents no
// difficulty. As a consequence, however, this means that no two neighbours
// may be eating simultaneously, since there are five philosophers and five forks.
//
// This is a simple implementation of Dijkstra's solution to the "Dining
// Philosophers" dilemma.

type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

var philosophers = []Philosopher{
	{name: "Plato", leftFork: 4, rightFork: 0},
	{name: "Socrates", leftFork: 0, rightFork: 1},
	{name: "Nietzsche", leftFork: 1, rightFork: 2},
	{name: "Aristotle", leftFork: 2, rightFork: 3},
	{name: "Stoic", leftFork: 3, rightFork: 4},
}

var numberOfMeals = 3
var timeOfEating = 1 * time.Second
var timeOfThinking = 3 * time.Second
var sleepTime = 1 * time.Second

var ordersMutex = &sync.Mutex{}
var ordersFinished = []string{}

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func main() {

	fmt.Println("Dining philosophers problem.")
	fmt.Println()
	fmt.Println("The table is empty, the philosophers have not came yet...")

	dine()

	fmt.Println("The table is empty, everyone ate...")

	fmt.Println(ordersFinished)
}

func dine() {

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(philosophers))

	seatsWaitGroup := &sync.WaitGroup{}
	seatsWaitGroup.Add(len(philosophers))

	var forks = make(map[int]*sync.Mutex)

	// Creating mutexes for each fork in a map
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	for i := 0; i < len(philosophers); i++ {
		go startEating(philosophers[i], waitGroup, seatsWaitGroup, forks)
	}

	waitGroup.Wait()

}

func startEating(
	philosopher Philosopher,
	waitGroup *sync.WaitGroup,
	seatsWaitGroup *sync.WaitGroup,
	forks map[int]*sync.Mutex) {

	defer waitGroup.Done()

	fmt.Printf("%s is seated at the table! \n", philosopher.name)
	seatsWaitGroup.Done()
	// Waiting for all philosophers to seat at the table
	seatsWaitGroup.Wait()

	for i := numberOfMeals; i > 0; i-- {

		// Trying to access both forks
		// Here a little special logic is needed so no two phisolophers
		// get stuck trying to acquire the same fork at the same time we
		// will always try to get the fork with lower number first
		forks[Min(philosopher.leftFork, philosopher.rightFork)].Lock()
		fmt.Printf("\t %s takes the first fork...\n", philosopher.name)

		forks[Max(philosopher.leftFork, philosopher.rightFork)].Lock()
		fmt.Printf("\t %s takes the second fork...\n", philosopher.name)

		fmt.Printf("\t %s has both forks! Ready to eat...\n", philosopher.name)

		time.Sleep(timeOfEating)

		fmt.Printf("\t %s is thinking...\n", philosopher.name)
		time.Sleep(timeOfThinking)

		forks[Min(philosopher.leftFork, philosopher.rightFork)].Unlock()
		forks[Max(philosopher.leftFork, philosopher.rightFork)].Unlock()

		fmt.Printf("\t %s stopped using the forks\n", philosopher.name)

	}

	fmt.Printf("\t %s is satisfied...\n", philosopher.name)
	fmt.Printf("\t %s left the table ... \n", philosopher.name)

	ordersMutex.Lock()
	ordersFinished = append(ordersFinished, philosopher.name)
	ordersMutex.Unlock()

}
