package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strconv"
	"time"
)

type Data struct {
	CollectionName string `bson:"collectionName" json:"collectionName"`
	Message        string `bson:"message" json:"message"`
	Date           string `bson:"date" json:"date"`
}

var ourData Data

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}
func main() {
	var (
		rabbitlogin, _   = os.LookupEnv("RABBITLOGIN")
		rabbitPass, _    = os.LookupEnv("RABBITPASS")
		rabbitHost, _    = os.LookupEnv("RABBITHOST")
		dbUser, _        = os.LookupEnv("DBUSER")
		dbPass, _        = os.LookupEnv("DBPASS")
		dbHost, _        = os.LookupEnv("DBHOST")
		dbName, _        = os.LookupEnv("DBNAME")
		collectionOne, _ = os.LookupEnv("COLLECTIONONE")
		collectionSec, _ = os.LookupEnv("COLLECTIONSEC")
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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			err := json.Unmarshal(d.Body, &ourData)
			//Create client options
			clientOptions := options.Client().ApplyURI("mongodb://" + dbUser + ":" + dbPass + "@" + dbHost)
			// Create connect to MongoDB
			client, err := mongo.Connect(context.TODO(), clientOptions)
			if err != nil {
				log.Fatal(err)
			}
			// Check the connection
			err = client.Ping(context.TODO(), nil)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Connected to MongoDB!")
			var ourCollection string

			if ourData.Date == "" {
				ourData.Date = currentDate()
			}

			if ourData.CollectionName == "testCollection#1" {
				ourCollection = collectionOne
			}
			if ourData.CollectionName == "testCollection#2" {
				ourCollection = collectionSec
			}
			collection := client.Database(dbName).Collection(ourCollection)
			newMsgG := Data{ourData.CollectionName, ourData.Message, ourData.Date}

			insertResult, err := collection.InsertOne(context.TODO(), newMsgG)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Inserted a single document: ", insertResult.InsertedID)
			if err != nil {
				fmt.Printf("Unmarshaling error: %s, %v", ourData, err)
				return

			}

			log.Printf("Received a message: %s", ourData)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
func currentDate() string {
	d := time.Now().Unix()
	return fmt.Sprintf(strconv.FormatInt(d, 10))
}
