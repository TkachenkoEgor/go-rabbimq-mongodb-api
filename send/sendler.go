package main

import (
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
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
}
