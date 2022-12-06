package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var router *mux.Router

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func CreateRouter() {
	router = mux.NewRouter()
}
func InitializeRoute() {
	router.HandleFunc("/mydate", requestHandler)
	router.HandleFunc("/login", loginHandler)

	serv := &http.Server{
		Handler:      router,
		Addr:         "localhost:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("server started")

	log.Fatal(serv.ListenAndServe())

}

func main() {
	CreateRouter()
	InitializeRoute()
}

func requestHandler(w http.ResponseWriter, r *http.Request) {

	errVal := validateToken(w, r)
	if errVal == nil {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{}

		ctx := context.Background()

		data := map[string]interface{}{}

		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			fmt.Println(err.Error())
		}

		switch r.Method {

		case "GET":
			response, err = getData(ctx, r)

		}

		if err != nil {
			response = map[string]interface{}{"error": err.Error()}
		}

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")

		if err := enc.Encode(response); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func getData(ctx context.Context, r *http.Request) (map[string]interface{}, error) {
	timeSince := r.URL.Query().Get("firstD")
	timeUntil := r.URL.Query().Get("secondD")
	collName := r.URL.Query().Get("collectionName")
	var (
		dbUser, _        = os.LookupEnv("DBUSER")
		dbPass, _        = os.LookupEnv("DBPASS")
		dbHost, _        = os.LookupEnv("DBHOST")
		dbName, _        = os.LookupEnv("DBNAME")
		collectionOne, _ = os.LookupEnv("COLLECTIONONE")
		collectionSec, _ = os.LookupEnv("COLLECTIONSEC")
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+dbUser+":"+dbPass+"@"+dbHost))

	if err != nil {
		fmt.Println(err.Error())
	}
	var ourCollection string

	if collName == "testCollection2" {
		ourCollection = collectionSec
	}
	if collName == "testCollection1" {
		ourCollection = collectionOne
	}
	collection := client.Database(dbName).Collection(ourCollection)

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"date", bson.D{{"$gte", timeSince}}}},
			bson.D{{"date", bson.D{{"$lte", timeUntil}}}},
		}}}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var products []bson.M
	var myDate []string
	for cur.Next(ctx) {

		var product bson.M

		if err = cur.Decode(&product); err != nil {
			return nil, err
		}

		products = append(products, product)

	}

	for _, mapa := range products {
		for index, value := range mapa {
			if index == "date" {
				switch vType := value.(type) {
				case string:
					strTime, _ := strconv.Atoi(vType)
					unxTime := int64(strTime)
					normTime := time.Unix(unxTime, 0).Format("02.01.2006")
					myDate = append(myDate, normTime)

				}
			}
		}
	}
	result := make(map[string]int)
	for _, value := range myDate {
		result[value]++
	}
	var sumDate float64
	for range myDate {
		sumDate = sumDate + 1
	}
	sumDays := float64(len(result))
	average := fmt.Sprintf("%.2f", sumDate/sumDays)
	res := map[string]interface{}{
		"data":          result,
		"average value": average,
	}
	return res, nil
}
