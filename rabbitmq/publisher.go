package rmq

import (
	"log"
	"simple-crud-rnd/config"

	"github.com/streadway/amqp"
)

type PublishMessage struct {
	Key     string
	Payload string
}

func createRabbitMQConnection(cfg *config.Config) (*RMQConfig, error) {
	return NewRabbitMQConfig(cfg.RabbitMQ.Username, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)
}

func SendMessage(cfg *config.Config, pm *PublishMessage) error {
	c, err := createRabbitMQConnection(cfg)
	if err != nil {
		return err
	}
	defer c.Close()

	err = c.Channel.Publish(
		"",
		cfg.RabbitMQ.Queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(pm.Payload),
			Expiration:  "60000",
			MessageId:   pm.Key,
		},
	)
	if err != nil {
		log.Printf("Error publish message: %v", err)
	}
	log.Printf("Message sent to queue %s: %s\n", cfg.RabbitMQ.Queue, pm.Payload)

	return nil
}
