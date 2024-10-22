package rabbitmq

import (
	"encoding/json"
	"log"
	"simple-crud-rnd/helpers"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

// RabbitMQRPC is a function to publish on rabbitmq
/**
 * @author Mahendra Dwi Purwanto
 *
 * @param connection *RabbitMQConnection, body interface{}, fReply bool
 *
 * @return responseStatus interface{}, responseMessage interface{}, responseData interface{}
 */
func RabbitMQRPC(connection *RabbitMQConnection, body interface{}, fReply bool) (responseStatus interface{}, responseMessage interface{}, responseData interface{}) {
	// Construct the URL for the RabbitMQ connection.
	url := connection.Host + ":" + connection.Port + connection.VirtualHost

	// Log the URL and the queue name.
	log.Println("[AMQP] " + url + " | " + connection.QueueName)

	// Establish a connection to RabbitMQ.
	dial, err := amqp.Dial("amqp://" + connection.Username + ":" + connection.Password + "@" + url)
	if err != nil {
		helpers.HandleError("failed to connect rabbitmq", err)
	}
	defer func(dial *amqp.Connection) {
		err := dial.Close()
		if err != nil {

		}
	}(dial)

	// Create a channel for communication with RabbitMQ.
	channel, err := dial.Channel()
	if err != nil {
		helpers.HandleError("failed to open a channel in rabbitmq", err)
	}
	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {

		}
	}(channel)

	// Generate a correlation ID and a unique queue name for this RPC request.
	corrId := helpers.RandomByte(8)
	quenameRd := "go-skeleton/" + strconv.Itoa(int(time.Now().Unix())) + corrId

	// Declare a unique queue for this RPC request.
	queue, err := channel.QueueDeclare(
		quenameRd,
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		helpers.HandleError("failed to declare a queue in rabbitmq", err)
	}
	defer func(channel *amqp.Channel, name string, ifUnused, ifEmpty, noWait bool) {
		_, err := channel.QueueDelete(name, ifUnused, ifEmpty, noWait)
		if err != nil {

		}
	}(channel, quenameRd, false, false, true)

	// Create a consumer to receive messages from the unique queue.
	message, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		helpers.HandleError("failed to register a consumer in rabbitmq", err)
	}

	// Declare the main queue if it doesn't exist.
	_, err = channel.QueueDeclare(
		connection.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		helpers.HandleError("failed to declare a queue in rabbitmq", err)
	}

	// Set up whether to wait for a response or not.
	waitResponse := true
	if fReply == false {
		waitResponse = false
		queue.Name = ""
	}

	// Publish the message to the main queue.
	err = channel.Publish(
		"",
		connection.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       queue.Name,
			Body:          []byte(helpers.JSONEncode(body)),
			Expiration:    "60000",
		},
	)
	log.Println("publish :", helpers.JSONEncode(body))
	if err != nil {
		helpers.HandleError("failed to publish a message to rabbitmq", err)
	}

	// Set a timeout for the RabbitMQ response.
	rabbitTimeout := time.After(10 * time.Second)

	// Wait for a response or timeout.
	for waitResponse {
		select {
		case <-rabbitTimeout:
			responseStatus = "500"
			responseMessage = "RPC timeout " + connection.QueueName
			waitResponse = false
		case data := <-message:
			log.Println("reply : ", string(data.Body))
			log.Println("")

			// Check if the received response corresponds to the request.
			if corrId == data.CorrelationId {
				if helpers.JSONEncode(body) == string(data.Body) {
					responseStatus = "error"
					responseMessage = "the rpc server did not respond"
					responseData = nil
				} else {
					// Parse the response and set the response variables.
					var response map[string]interface{}
					err := json.Unmarshal([]byte(string(data.Body)), &response)
					if err != nil {
						return nil, nil, nil
					}
					responseStatus = response["Status"]
					responseMessage = response["StatusMessage"]
					responseData = response
				}
				waitResponse = false
			}
		}
	}
	return
}
