package rmq

import (
	"fmt"
	"log"
	"simple-crud-rnd/config"
	"simple-crud-rnd/rabbitmq/consumer"
)

func NewRabbitMQConsumer(cfg *config.Config) (*RMQConfig, error) {
	return NewRabbitMQConfig(cfg.RabbitMQ.Username, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)
}

func RunConsumer(cfg *config.Config, c *RMQConfig) {
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalln("Error opening database")
	}

	msgs, err := c.Channel.Consume(
		cfg.RabbitMQ.Queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Gagal consume channel rabbitmq: %v", err)
	}

	log.Println("Waiting for messages from queue:", cfg.RabbitMQ.Queue)
	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			fmt.Printf("Received new message: %s\n", msg.Body)
			consumer.ConsumerRouter(&msg, db)
		}
	}()

	fmt.Println("Success to connect RabbitMQ.......")
	<-forever
}
