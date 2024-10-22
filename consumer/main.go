package main

import (
	"log"
	"os"
	"simple-crud-rnd/rabbitmq"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error load .env: %v", err)
	}
	rabbitmq.RabbitMQ(os.Getenv("RQ_QUEUE"))
}
