package main

import (
	"ibdatabase/api"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// set application config
	var app api.Application

	err := godotenv.Load()
	if err != nil {
		log.Println("no local env: ", err)
	}

	app = api.Application{}

	app.StartApp()
}
