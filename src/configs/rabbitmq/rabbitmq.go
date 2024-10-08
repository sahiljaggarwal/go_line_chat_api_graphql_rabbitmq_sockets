// src/configs/rabbitmq.go
package rabbitmq

import (
	"line/src/configs/env"
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection

func Connect() {
	var err error
	conn, err = amqp.Dial(env.RABBITMQ_URL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"chat", // name of the queue
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}
	log.Print("RabbitMQ connected and queue declared")
}

func GetConnection() *amqp.Connection {
	return conn
}

func CloseConnection() {
	if conn != nil {
		conn.Close()
	}
}
