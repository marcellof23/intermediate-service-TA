package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	// Create a channel to listen for signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	// Run an infinite loop
	for {
		select {
		case <-signalCh:
			// Ctrl+C signal received, break out of the loop
			fmt.Println("Received Ctrl+C signal. Breaking out of the loop.")
			return
		default:
			// Your code here that runs in the infinite loop
			fmt.Println("Running the infinite loop...")
		}
	}
}
