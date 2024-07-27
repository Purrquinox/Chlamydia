package main

import (
	"Chlamydia/api"
	"Chlamydia/state"

	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Setup
	state.Setup()
	api.StartAPI()

	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-sigs
	fmt.Printf("Received signal: %s, shutting down gracefully...\n", sig)
}
