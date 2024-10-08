package main

import (
	"line/src/app"
	"line/src/common/shutdown"
	"line/src/configs/env"
	"log"
	"time"
)

func main() {
	app := app.SetupApp()
	port := env.PORT
	if port == "" {
		port = "3000"
	}
	serve := ":" + port
	// log.Fatal(app.Listen(serve))

	go func() {
		if err := app.Listen(serve); err != nil {
			log.Printf("Error starting server: %v\n", err)
		}
	}()

	log.Printf("Server is running on port %s", port)

	shutdown.GracefulShutdown(app, 5*time.Second)
}
