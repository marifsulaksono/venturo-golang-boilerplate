package main

import (
	"log"
	"os"
	"os/signal"
	"simple-crud-rnd/config"
	rmq "simple-crud-rnd/rabbitmq"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error loading configs")
	}

	consumer, err := rmq.NewRabbitMQConsumer(cfg)
	if err != nil {
		log.Fatalf("Error starting RabbitMQ consumer: %v", err)
	}
	defer consumer.Close()

	// create a channel to listen for interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// run the consumer in a separate goroutine
	go func() {
		rmq.RunConsumer(cfg, consumer)
	}()

	// block until signal to stop is received
	sig := <-signalChan
	time.Sleep(1 * time.Second)
	log.Printf("Received signal: %s, shutting down gracefully...", sig)
}
