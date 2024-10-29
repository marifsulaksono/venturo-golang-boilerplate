package rmq

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type RMQConfig struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      string
}

func NewRabbitMQConfig(username, password, host, port string) (*RMQConfig, error) {
	queue := os.Getenv("RQ_QUEUE")
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Error declare queue: %v", err)
	}

	fmt.Println("declared queue name:", queue)

	return &RMQConfig{
		Connection: conn,
		Channel:    ch,
		Queue:      queue,
	}, nil
}

func (conf *RMQConfig) Close() {
	if conf.Channel != nil {
		conf.Channel.Close()
	}
	if conf.Connection != nil {
		conf.Connection.Close()
	}
}
