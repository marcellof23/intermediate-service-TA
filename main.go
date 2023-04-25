package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

//func main() {
//	// Create a channel to listen for signals
//	signalCh := make(chan os.Signal, 1)
//	signal.Notify(signalCh, os.Interrupt)
//
//	// Run an infinite loop
//	for {
//		select {
//		case <-signalCh:
//			// Ctrl+C signal received, break out of the loop
//			fmt.Println("Received Ctrl+C signal. Breaking out of the loop.")
//			return
//		default:
//			// Your code here that runs in the infinite loop
//			fmt.Println("Running the infinite loop...")
//			time.Sleep(3000 * time.Millisecond)
//		}
//	}
//
//	fmt.Println("AAAA")
//}

func handler(signal os.Signal) {
	if signal == syscall.SIGTERM {
		fmt.Println("Got kill signal. ")
		fmt.Println("Program will terminate now.")
		os.Exit(0)
	} else if signal == syscall.SIGINT {
		fmt.Println("Got CTRL+C signal")
		fmt.Println("Closing.")
		os.Exit(0)
	} else {
		fmt.Println("Ignoring signal: ", signal)
	}
}

func tes() chan int {
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)
	exitchnl := make(chan int)
	go func() {
		for {
			select {
			case <-sigchnl:
				s := <-sigchnl
				handler(s)
			default:
				fmt.Println("AA")
			}

		}
	}()
	return exitchnl
}

func main() {
	exitcode := <-tes()
	fmt.Println("Ignoring signal: ")
	os.Exit(exitcode)
}
