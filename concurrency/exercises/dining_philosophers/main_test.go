package main

import (
	"testing"
	"time"
)

func TestDine(t *testing.T) {
	timeOfEating = 0 * time.Second
	sleepTime = 0 * time.Second
	timeOfThinking = 0 * time.Second

	for i := 0; i < 10; i++ {
		ordersFinished = []string{}
		dine()

		if len(ordersFinished) != len(philosophers) {
			t.Errorf(
				"Expected %d but got %d...",
				len(ordersFinished),
				len(philosophers),
			)
		}

	}
}
