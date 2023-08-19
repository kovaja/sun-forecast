package server

import (
	"errors"
	"kovaja/sun-forecast/api"
	"kovaja/sun-forecast/forecast"
	"kovaja/sun-forecast/logger"
	"kovaja/sun-forecast/utils"
	"kovaja/sun-forecast/weather"
	"net/http"
)

const DEFAULT_API_PATH = "/api/"

func defaultPathHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, DEFAULT_API_PATH, http.StatusMovedPermanently)
}

func defaultApiHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log("Default api called")
	api.SendResponse(w, nil, nil)
}

func currentWeatherHandler(w http.ResponseWriter, r *http.Request) {
	data, err := weather.GetCurrentWeather()
	api.SendResponse(w, data, err)
}

func weatherForcastHandler(w http.ResponseWriter, r *http.Request) {
	data, err := weather.GetForecast()
	api.SendResponse(w, data, err)
}

func forecastHandler(w http.ResponseWriter, r *http.Request) {
	data, err := forecast.ReadForecastsFromDb()
	api.SendResponse(w, data, err)
}

func consumeForecastHandler(w http.ResponseWriter, r *http.Request) {
	err := forecast.ConsumeForecasts()
	api.SendResponse(w, nil, err)
}

func updateForecastHandler(w http.ResponseWriter, r *http.Request) {
	// logger.Log("Update forecast api called %s | %s | %s | %s | %v", r.Method, r.URL, r.RequestURI, r.RemoteAddr, r.TLS)
	logger.Log("Update forecast api called %v", r)

	var err error
	if r.Method == http.MethodPost {
		updates, err := forecast.UpdateForecasts(r)
		api.SendResponse(w, updates, err)
	} else {
		err = errors.New("Method not allowed")
		api.SendError(w, err, http.StatusMethodNotAllowed)
	}

}

func InitializeServer() error {
	port, err := utils.GetEnvVariable("PORT")

	if err != nil {
		return utils.CustomError("Failed to load port env variable", err)
	}

	http.HandleFunc(DEFAULT_API_PATH+"weather/", currentWeatherHandler)
	http.HandleFunc(DEFAULT_API_PATH+"weather/forecast/", weatherForcastHandler)
	http.HandleFunc(DEFAULT_API_PATH+"forecast/", forecastHandler)
	http.HandleFunc(DEFAULT_API_PATH+"forecast/consume/", consumeForecastHandler)
	http.HandleFunc(DEFAULT_API_PATH+"forecast/update/", updateForecastHandler)
	http.HandleFunc(DEFAULT_API_PATH, defaultApiHandler)
	http.HandleFunc("/", defaultPathHandler)

	logger.Log("Server will listen on port %s", (":" + port))
	return http.ListenAndServe(":"+port, nil)
}
