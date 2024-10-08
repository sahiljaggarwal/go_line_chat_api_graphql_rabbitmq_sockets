// src/main.go
package main

import (
	"line/src/app"
	"line/src/configs/env"
	"log"
)

func main() {
	app := app.SetupApp()
	port := env.PORT
	if port == "" {
		port = "3000"
	}
	serve := ":" + port
	log.Fatal(app.Listen(serve))
}
