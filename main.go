package main

import (
	"encoding/json"
	"fmt"
	"kovaja/sun-forecast/weather"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

const DEFAULT_API_PATH = "/api/"

type Response map[string]any
type ErrorResponse struct {
	Error string `json:"error"`
}

func getResponse(data any, err error) (Response, error) {
	response := make(Response)

	response["date"] = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	response["data"] = data

	return response, err
}

func sendJson(w http.ResponseWriter, d Response, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorRespone := ErrorResponse{Error: err.Error()}
		errorResponseJson, err := json.Marshal(errorRespone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(errorResponseJson)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(d)
}

func defaultPathHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, DEFAULT_API_PATH, http.StatusMovedPermanently)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	response, err := getResponse(nil, nil)
	sendJson(w, response, err)
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	data, err := weather.GetFourDays()
	response, err := getResponse(data, err)
	sendJson(w, response, err)
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
	http.HandleFunc(DEFAULT_API_PATH, apiHandler)
	http.HandleFunc("/", defaultPathHandler)

	checkError(http.ListenAndServe(":8080", nil))
}
