package main

import (
	"log"

	"simple-crud-rnd/config"
	"simple-crud-rnd/routes"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading configs")
	}

	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Error opening database")
	}

	e := routes.NewHTTPServer(cfg, db)
	e.RunHTTPServer()
}
