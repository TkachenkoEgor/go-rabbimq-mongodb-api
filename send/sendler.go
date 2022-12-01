package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	var (
		rabbitlogin, _ = os.LookupEnv("RABBITLOGIN")
		rabbitPass, _  = os.LookupEnv("RABBITPASS")
		rabbitHost, _  = os.LookupEnv("RABBITHOST")
	)
	conn, err := amqp.Dial("amqp://" + rabbitlogin + ":" + rabbitPass + "@" + rabbitHost)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msgBody := `{"collectionName":"testCollection#2",
			"message":"test message2",
			"date":"66329999999999999"}`
	fmt.Println(msgBody)
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msgBody),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", msgBody)

}
