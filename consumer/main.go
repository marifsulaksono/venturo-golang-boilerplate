package main

import (
	"log"
	"simple-crud-rnd/config"
	"simple-crud-rnd/rabbitmq"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error loading configs")
	}

	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalln("Error opening database")
	}

	rabbitmq.RabbitMQ(cfg.RabbitMQ.Queue, db)
}
