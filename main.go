package main

import (
	"log"

	"simple-crud-rnd/config"
	"simple-crud-rnd/routes"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error loading configs:", err)
	}

	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalln("Error connect database:", err)
	}

	rds, err := config.InitRedisClient(cfg)
	if err != nil {
		log.Fatalln("Error opening database:", err)
	}

	e := routes.NewHTTPServer(cfg, db, rds)
	e.RunHTTPServer()
}
