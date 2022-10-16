package main

import (
	"enterprise.sidooh/api"
	"enterprise.sidooh/pkg/db"
	"enterprise.sidooh/pkg/logger"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.Init()
	db.Init()

	app := api.Server()

	go func() {
		log.Fatal(app.Listen(":8006"))
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()
}
