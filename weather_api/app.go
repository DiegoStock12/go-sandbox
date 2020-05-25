package weather_api

import (
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

// Provide an api endpoint that allows for filtering the
// results per city and possibly per coordinates

// The api returns JSON by default

// Coordinate object
type coord struct{
	Latitude float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// Temp represents the max and min temp
type temp struct{
	Max float32 `json:"max"`
	Min float32 `json:"min"`
}

// Single day pred
type dayPrediction struct {
	Temperature temp `json:"temperature"`
	Rain float32 `json:"rain"`
	Description string `json:"description"`
}

// WeatherPrediction encapsulates the prediction for that city
// in the following three days
type WeatherPrediction struct{
	City string `json:"city"`
	Coordinates coord `json:"coordinates"`
	Prediction []dayPrediction `json:"prediction"`
}

// Define the cache for keeping the weather prediction
// Might implement a MySQL or MongoDB in the future
// Even a GEO-Spatial DB for closeness search
var cache map[string]WeatherPrediction

// Mutex to allow reads
var mutex = &sync.Mutex{}

// Returns prediction per city
func getPredictionByName(w http.ResponseWriter, r *http.Request){

}


func main(){
	// Initialize the router
	router := mux.NewRouter()

	// go on


}
