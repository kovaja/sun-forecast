package server

import (
	"database/sql"
	"errors"
	app "kovaja/sun-forecast/business"
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

type Server struct {
	appControllers app.AppControllers
}

func (s Server) defaultPathHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log("Serving index.html %s", r.URL)
	http.ServeFile(w, r, "static/index.html")
}

func (s Server) defaultApiHandler(r *http.Request) (any, error) {
	return nil, nil
}

func currentWeatherHandler(r *http.Request) (any, error) {
	return weather.GetCurrentWeather()
}

func weatherForecastHandler(r *http.Request) (any, error) {
	return weather.GetForecast()
}

func (s Server) forecastHandler(r *http.Request) (any, error) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if toStr == "" {
		return nil, ErrMissingParamTo
	}

	if fromStr == "" {
		return nil, ErrMissingParamFrom
	}

	return s.appControllers.ForecastCtl.GetForecasts(fromStr, toStr)
}

func (s Server) consumeForecastHandler(r *http.Request) (any, error) {
	err := s.appControllers.ForecastCtl.ConsumeForecasts()
	return nil, err
}

func (s Server) updateForecastHandler(r *http.Request) (any, error) {
	if r.Method == http.MethodPost {
		return s.appControllers.ForecastCtl.UpdateForecasts(r)
	}

	return nil, ErrMethodNotAllowed
}

func (s Server) eventHandler(r *http.Request) (any, error) {
	typeStr := r.URL.Query().Get("type")
	limitStr := r.URL.Query().Get("limit")
	return s.appControllers.EventCtl.ReadEvents(typeStr, limitStr)
}

func InitializeServer(db *sql.DB) error {
	port, err := utils.GetEnvVariable("PORT")

	if err != nil {
		return utils.CustomError("Failed to load port env variable", err)
	}

	if utils.IsDev() && port == "" {
		port = "8080"
	}

	server := Server{
		appControllers: app.InitializeApp(db),
	}

	var routes map[string]ApiHandler = map[string]ApiHandler{
		"forecast":         server.forecastHandler,
		"weather":          currentWeatherHandler,
		"weather/forecast": weatherForecastHandler,
		"forecast/consume": server.consumeForecastHandler,
		"forecast/update":  server.updateForecastHandler,
		"event":            server.eventHandler,
		"":                 server.defaultApiHandler,
	}

	for path, handler := range routes {
		logger.Log("Register handler %s/", path)
		http.HandleFunc(DEFAULT_API_PATH+path+"/", logRequest(handler))
	}

	http.HandleFunc("/", server.defaultPathHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	logger.Log("Server will listen on port %s", (":" + port))
	return http.ListenAndServe(":"+port, nil)
}
