package solcast

import (
	"errors"
	"kovaja/sun-forecast/events"
	"time"
)

const tsLayout = "2006-01-02T15:04:05.0000000Z"

func constructForcast(apiForecast SolcastApiForecast) (*Forecast, error) {
	parsedTime, err := time.Parse(tsLayout, apiForecast.PeriodEnd)
	if err != nil {
		return nil, err
	}

	forecast := Forecast{
		Id:        -1,
		PeriodEnd: parsedTime,
		Value:     apiForecast.PvEstimate * 1000, // convert 1.56kw na 1560w
		Actual:    nil,
	}

	return &forecast, nil
}

func ReadForecastsFromApi() (*ForecastResponse, error) {
	canCall, remainingCalls := CanCall()

	if !canCall {
		return nil, errors.New("Cannot call solcast api (Too many calls today)")
	}

	body, err := fetchForecasts()
	if err != nil {
		return nil, err
	}

	var forcasts []Forecast
	for _, apiForecast := range body.Forecasts {
		forecast, err := constructForcast(apiForecast)

		if err != nil {
			events.LogEvent("Failed to construct forecast %v. With error: %v", apiForecast, err)
		} else {
			forcasts = append(forcasts, *forecast)
		}
	}

	response := ForecastResponse{
		RemainingCalls: remainingCalls,
		Forecasts:      forcasts,
	}

	return &response, nil
}
