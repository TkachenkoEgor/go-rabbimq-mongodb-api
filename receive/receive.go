package receive

import (
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

var ourData Data

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// loads values from .env into the system
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var (
	rabbitlogin, _         = os.LookupEnv("RABBITLOGIN")
	rabbitPass, _          = os.LookupEnv("RABBITPASS")
	rabbitHost, _          = os.LookupEnv("RABBITHOST")
	dbUser, _              = os.LookupEnv("DBUSER")
	dbPass, _              = os.LookupEnv("DBPASS")
	dbHost, _              = os.LookupEnv("DBHOST")
	dbName, _              = os.LookupEnv("DBNAME")
	collectionNewScans, _  = os.LookupEnv("COLLECTIONNEWSCANS")
	collectionNewClient, _ = os.LookupEnv("COLLECTIONEWCLIENT")
)

func main() {
	conn, err := amqp.Dial("amqp://" + rabbitlogin + ":" + rabbitPass + "@" + rabbitHost)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
}
