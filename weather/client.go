package weather

import (
	"errors"
	"fmt"
	"kovaja/sun-forecast/httpClient"
	"kovaja/sun-forecast/logger"
	"kovaja/sun-forecast/utils"
)

const LAT_VAR = "LAT"
const LON_VAR = "LON"
const API_KEY_VAR = "WEATHER_API_KEY"

type WeatherResponse struct {
	RemainingCalls int `json:"remainingCalls"`
	Data           any `json:"data"`
}

func getParams() (string, error) {
	apiKey, err := utils.GetEnvVariable(API_KEY_VAR)
	lat, err := utils.GetEnvVariable(LAT_VAR)
	lon, err := utils.GetEnvVariable(LON_VAR)

	params := fmt.Sprintf("units=metric&lat=%s&lon=%s&appid=%s", lat, lon, apiKey)

	return utils.ReturnStringResultOrError(params, err)
}

func getUrl(path string) (string, error) {
	// https://api.openweathermap.org/data/2.5/weather?lat={lat}&lon={lon}&appid={API key}
	params, err := getParams()
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/%s?%s", path, params)

	return utils.ReturnStringResultOrError(url, err)
}

func getWeatherData(path string) (*WeatherResponse, error) {
	canCall, remainingCalls := CanCall()
	if !canCall {
		return nil, errors.New("Cannot call weather api")
	}

	url, err := getUrl(path)
	if err != nil {
		return nil, errors.New("Failed to build weather api url")
	}

	var body interface{}
	err = httpClient.GetJson(url, &body)

	if err != nil {
		errorMsg := fmt.Sprintf("Failed to call %s, response: %v", path, body)
		logger.LogError(errorMsg, err)
		return nil, errors.New("Failed to call weather api")
	}

	r := &WeatherResponse{
		RemainingCalls: remainingCalls,
		Data:           body,
	}

	return r, nil
}

func GetCurrentWeather() (*WeatherResponse, error) {
	path := "weather"
	return getWeatherData(path)
}

func GetForecast() (*WeatherResponse, error) {
	path := "forecast"
	return getWeatherData(path)
}
