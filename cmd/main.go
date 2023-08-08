package main

import (
	"flag"
	"ibdatabase/api"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// set application config
	var app api.Application

	env := flag.String("env", "dev", "init")
	flag.Parse()

	if *env == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Println("no local env: ", err)
		}
	}

	app = api.Application{}

	app.StartApp()
}
