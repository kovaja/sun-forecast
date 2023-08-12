package main

import (
	"fmt"
	"kovaja/sun-forecast/api"
	"kovaja/sun-forecast/weather"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

const DEFAULT_API_PATH = "/api/"

func defaultPathHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, DEFAULT_API_PATH, http.StatusMovedPermanently)
}

func defaultApiHandler(w http.ResponseWriter, r *http.Request) {
	api.SendResponse(w, nil, nil)
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	data, err := weather.GetFourDays()
	api.SendResponse(w, data, err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
}

func main() {
	checkError(godotenv.Load())

	http.HandleFunc(DEFAULT_API_PATH+"weather/", weatherHandler)
	http.HandleFunc(DEFAULT_API_PATH, defaultApiHandler)
	http.HandleFunc("/", defaultPathHandler)

	checkError(http.ListenAndServe(":8080", nil))
}
