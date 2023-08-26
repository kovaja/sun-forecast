package server

import (
	"database/sql"
	"errors"
	"kovaja/sun-forecast/business/events"
	"kovaja/sun-forecast/business/forecast"
	"kovaja/sun-forecast/business/weather"
	"kovaja/sun-forecast/utils"
	"kovaja/sun-forecast/utils/logger"
	"net/http"
)

const DEFAULT_API_PATH = "/api/"

var (
	ErrMissingParamTo   = errors.New("Missing required parameter to")
	ErrMissingParamFrom = errors.New("Missing required parameter from")
	ErrMethodNotAllowed = errors.New("Method not allowed")
)

type ApiHandler func(r *http.Request) (any, error)

var forecastController forecast.ForecastController
var eventController events.EventController

func defaultPathHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log("Serving index.html %s", r.URL)
	http.ServeFile(w, r, "static/index.html")
}

func defaultApiHandler(r *http.Request) (any, error) {
	return nil, nil
}

func currentWeatherHandler(r *http.Request) (any, error) {
	return weather.GetCurrentWeather()
}

func weatherForecastHandler(r *http.Request) (any, error) {
	return weather.GetForecast()
}

func forecastHandler(r *http.Request) (any, error) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if toStr == "" {
		return nil, ErrMissingParamTo
	}

	if fromStr == "" {
		return nil, ErrMissingParamFrom
	}

	return forecastController.GetForecasts(fromStr, toStr)
}

func consumeForecastHandler(r *http.Request) (any, error) {
	err := forecastController.ConsumeForecasts()
	return nil, err
}

func updateForecastHandler(r *http.Request) (any, error) {
	if r.Method == http.MethodPost {
		return forecastController.UpdateForecasts(r)
	}

	return nil, ErrMethodNotAllowed
}

func eventHandler(r *http.Request) (any, error) {
	return events.ReadEvents()
}

var routes map[string]ApiHandler = map[string]ApiHandler{
	"forecast":         forecastHandler,
	"weather":          currentWeatherHandler,
	"weather/forecast": weatherForecastHandler,
	"forecast/consume": consumeForecastHandler,
	"forecast/update":  updateForecastHandler,
	"event":            eventHandler,
	"":                 defaultApiHandler,
}

func InitializeServer(db *sql.DB) error {
	port, err := utils.GetEnvVariable("PORT")

	if err != nil {
		return utils.CustomError("Failed to load port env variable", err)
	}

	eventController = events.InitializeController(db)
	forecastController = forecast.InitializeController(db, eventController)

	for path, handler := range routes {
		logger.Log("Register handler %s/", path)
		http.HandleFunc(DEFAULT_API_PATH+path+"/", logRequest(handler))
	}

	http.HandleFunc("/", defaultPathHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	logger.Log("Server will listen on port %s", (":" + port))
	return http.ListenAndServe(":"+port, nil)
}
