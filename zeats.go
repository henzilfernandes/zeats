package main

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"zeats/wire"
)

func main() {

	fmt.Println(`
     _____             _
    |__  / ___   __ _ | |_  ___
      / / / _ \ / _  || __|/ __|
     / /_|  __/| (_| || |_ \__ \
    /____|\___| \__,_| \__||___/
	`)

	// Signal listener
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	wire.Start()

	// Listens to Signal
	go func() {
		sig := <-sigs
		fmt.Println("Server shutting down. Signal:", sig)
		done <- true
	}()

	// Close on shutdown signal
	<-done
	close(sigs)
	close(done)
	wire.Stop()
}

