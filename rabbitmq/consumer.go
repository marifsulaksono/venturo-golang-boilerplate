package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/rabbitmq/rmqauto"
	"simple-crud-rnd/structs"

	"github.com/streadway/amqp"
)

func RabbitMQ(queueName string) {
	log.Println("connect to rabbit mq")

	rbMq := rmqauto.CreateRqPubConsumer()
	rbMq.SetReadQueue(queueName)
	_, err := rbMq.StartConnection(os.Getenv("RQ_USERNAME"), os.Getenv("RQ_PASSWORD"), os.Getenv("RQ_HOST"), os.Getenv("RQ_PORT"),
		os.Getenv("RQ_VHOST"))
	if err != nil {
		helpers.HandleError("failed to connect rabbitmq", err)
	}
	defer rbMq.Stop()

	err = createQueue(rbMq, queueName)
	if err != nil {
		helpers.HandleError(fmt.Sprintf("failed to declare a queue %s in rabbitmq", queueName), err)
	}

	rbMq.ConsumeMessage()
	message := rbMq.GetMessageChanel(queueName)
	forever := make(chan bool)
	defer close(forever)

	deliveredMsg := make(chan amqp.Delivery)
	defer close(deliveredMsg)

	processMessages(rbMq, deliveredMsg)

	go func() {
		for data := range message {
			data.Ack(false)
			deliveredMsg <- data
		}
	}()
	log.Println("module is ready now")
	fmt.Println("")
	<-forever
}

func processMessages(rbMq rmqauto.IRqAutoConnect, msgCh <-chan amqp.Delivery) {
	for i := 0; i < 300; i++ {
		go processMessage(rbMq, msgCh)
	}
}

func processMessage(rbMq rmqauto.IRqAutoConnect, msgCh <-chan amqp.Delivery) {
	for data := range msgCh {
		processData(rbMq, data)
	}
}

func processData(rbMq rmqauto.IRqAutoConnect, data amqp.Delivery) {
	defer func() { recover() }()

	var request structs.Request
	err1 := json.Unmarshal(data.Body, &request)
	if err1 != nil {
		log.Println("ERROR :", err1)
	}
	log.Println("data.ReplyTo : ", data.ReplyTo)
	response := HandleRequest(request)
	log.Println("Response :", response)

	if data.ReplyTo == "" {
		return
	}

	err := rbMq.GetRqChannel().Publish(
		"",
		data.ReplyTo,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: data.CorrelationId,
			Body:          []byte(response),
			Expiration:    "60000",
		},
	)
	if err != nil {
		helpers.HandleError("failed to publish a message to rabbitmq", err)
	}
	log.Println("Success publish message reply")
	log.Println("")
}

func createQueue(rbMq rmqauto.IRqAutoPubConsumer, queueName string) (err error) {
	_, err = rbMq.GetRqChannel().QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	return
}
