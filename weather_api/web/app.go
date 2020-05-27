package main

import (
	"encoding/json"
	"fmt"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Provide an api endpoint that allows for filtering the
// results per city and possibly per coordinates

// The api returns JSON by default

// Coordinate object
type coord struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// Temp represents the max and min temp
type temp struct {
	Max float32 `json:"max"`
	Min float32 `json:"min"`
}

// Single day pred
type dayPrediction struct {
	Temperature temp   `json:"temperature"`
	Rain        int    `json:"rain"`
	Description string `json:"description"`
}

// WeatherPrediction encapsulates the prediction for that city
// in the following three days
type WeatherPrediction struct {
	City        string          `json:"city"`
	Coordinates coord           `json:"coordinates"`
	Prediction  []dayPrediction `json:"prediction"`
}

// Define the cache for keeping the weather prediction
// Might implement a MySQL or MongoDB in the future
// Even a GEO-Spatial DB for closeness search
var cache map[string]WeatherPrediction

// Mutex to allow reads
var mutex = &sync.Mutex{}

// Get the database name from the environment
var dbHost = os.Getenv("DB_NAME")
var connectionURI = fmt.Sprintf("mongodb://diego:passwd@%s:27017/?authSource=admin", dbHost)


// insert the cache objects example into the mongo db
func insertToDB() {
	// open a client and connect
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
	if err != nil{log.Fatal(err)}

	defer client.Disconnect(ctx)

	// insert
	collection := client.Database("test").Collection("weather")

	// Get all the values
	values := make([]interface{}, 0, len(cache))
	for _, v := range cache {
		values = append(values, v)
	}

	result, err := collection.InsertMany(ctx, values)
	if err != nil{
		log.Fatal(err)
	}

	log.Printf("IDs is %v", result)

}

// test getting the prediction from the mongodb db instead
func getPredictionDB(w http.ResponseWriter, r *http.Request) {

	// City name
	vars := mux.Vars(r)
	city := vars["city"]

	// open a client and connect
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))

	// get the collection
	collection := client.Database("test").Collection("weather")

	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)

	// result from the query
	var result WeatherPrediction

	// Get the result
	err = collection.FindOne(ctx, bson.M{"city": city}).Decode(&result)
	if err != nil {
		// Return an error
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Return the JSON for that
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// Returns prediction per city
func getPredictionByName(w http.ResponseWriter, r *http.Request) {
	// get the variables
	vars := mux.Vars(r)
	w.Header().Add("Content-Type", "application/json")

	// Get the name
	name := vars["city"]

	mutex.Lock()
	data := cache[strings.ToLower(name)]
	mutex.Unlock()

	// Get the Weather prediction and map it to JSON
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func generateData() {
	cities := []string{"chicago", "delft", "amsterdam"}

	for _, city := range cities {
		latitude := rand.Float32() * 90
		longitude := rand.Float32() * 90
		pred := WeatherPrediction{
			City: city,
			Coordinates: coord{
				Latitude:  latitude,
				Longitude: longitude,
			},
			Prediction: []dayPrediction{},
		}
		for i := 0; i < 3; i++ {
			maxTemp := float32(rand.Intn(26))
			minTemp := maxTemp - float32(rand.Intn(10))
			rain := rand.Intn(100)
			description := "quite good"
			// Append that to the object
			pred.Prediction = append(pred.Prediction, dayPrediction{
				Temperature: temp{Min: minTemp, Max: maxTemp},
				Rain:        rain,
				Description: description,
			})
		}

		// Append to the cache
		cache[city] = pred
	}

}

func main() {

	// Initialize the map
	cache = make(map[string]WeatherPrediction)

	log.Println(connectionURI)

	// generate some random data to put in
	generateData()

	// Insert it to the database
	insertToDB()


	fmt.Println(cache)

	// Initialize the router
	router := mux.NewRouter()
	router.HandleFunc("/weather/{city}", getPredictionByName).Methods("Get")
	router.HandleFunc("/weather/mongo/{city}", getPredictionDB).Methods("Get")


	log.Fatal(http.ListenAndServe(":8080", router))

}
