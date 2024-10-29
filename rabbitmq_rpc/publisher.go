package rabbitmq

import (
	"os"
	"simple-crud-rnd/structs"
	"time"

	"github.com/google/uuid"
)

type RabbitMQConnection struct {
	Host, Port, Username, Password, VirtualHost, QueueName string
}

// GlobalConn is a function to set global connection
/**
 * @author Mahendra Dwi Purwanto
 *
 * @param connection *RabbitMQConnection
 *
 * @return -
 */
func (connection *RabbitMQConnection) GlobalConn() {
	connection.Host = os.Getenv("RQ_HOST")
	connection.Port = os.Getenv("RQ_PORT")
	connection.Username = os.Getenv("RQ_USERNAME")
	connection.Password = os.Getenv("RQ_PASSWORD")
	connection.VirtualHost = os.Getenv("RQ_VHOST")
	connection.QueueName = os.Getenv("RQ_QUEUE")
}

func RequestCommand(route string, param string, data interface{}, id string, freplay bool) (interface{}, interface{}, interface{}) {
	connRabbitmq := RabbitMQConnection{}
	connRabbitmq.GlobalConn()
	return RabbitMQRPC(&connRabbitmq, structs.RabbitMQDefaultPayload{
		JobId: id,
		Route: route,
		Param: param,
		Data:  data,
	}, freplay)
}

func SendRabbitMQ(conn RabbitMQConnection, route string, param string, data interface{}, freplay bool) (interface{}, interface{}, interface{}) {
	return RabbitMQRPC(&conn, structs.MessagePayload{
		Id:         time.Now().UnixMicro(),
		Command:    route,
		Time:       time.Now().String(),
		ModuleId:   "pus-be-religious-manager",
		Properties: nil,
		Signature:  uuid.New().String(),
		Data:       data,
	}, freplay)
}
