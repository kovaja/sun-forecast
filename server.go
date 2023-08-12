package main

import (
	"kovaja/sun-forecast/api"
	"kovaja/sun-forecast/weather"
	"net/http"
)

const DEFAULT_API_PATH = "/api/"

func defaultPathHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, DEFAULT_API_PATH, http.StatusMovedPermanently)
}

func defaultApiHandler(w http.ResponseWriter, r *http.Request) {
	api.SendResponse(w, nil, nil)
}

func currentWeatherHandler(w http.ResponseWriter, r *http.Request) {
	data, err := weather.GetCurrentWeather()
	api.SendResponse(w, data, err)
}

func forceastHandler(w http.ResponseWriter, r *http.Request) {
	data, err := weather.GetForecast()
	api.SendResponse(w, data, err)
}

func InitializeServer() error {
	http.HandleFunc(DEFAULT_API_PATH+"weather/", currentWeatherHandler)
	http.HandleFunc(DEFAULT_API_PATH+"forecast/", forceastHandler)
	http.HandleFunc(DEFAULT_API_PATH, defaultApiHandler)
	http.HandleFunc("/", defaultPathHandler)

	return http.ListenAndServe(":8080", nil)
}
