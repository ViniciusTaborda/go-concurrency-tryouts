package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestPrintMessage(t *testing.T) {
	standardOutput := os.Stdout

	read, write, _ := os.Pipe()

	os.Stdout = write

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	textInput := "alpha"

	go UpdateMessage(textInput, &waitGroup)
	waitGroup.Wait()

	PrintMessage()

	_ = write.Close()

	result, _ := io.ReadAll(read)

	output := string(result)

	os.Stdout = standardOutput

	if !strings.Contains(output, textInput) {
		t.Errorf("Expected %s received %s.", textInput, output)
	}

}
