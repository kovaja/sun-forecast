package forecast

import "time"

type SolcastApiForecast struct {
	PvEstimate   float64 `json:"pv_estimate"`
	PvEstimate10 float64 `json:"pv_estimate10"`
	PVEstimate90 float64 `json:"pv_estimate90"`
	PeriodEnd    string  `json:"period_end"`
	Period       string  `json:"period"`
}

type SolcastApiForcastResponse struct {
	Forecasts []SolcastApiForecast `json:"forecasts"`
}

type Forecast struct {
	Id        int       `json:"id"`
	PeriodEnd time.Time `json:"periodEnd"`
	Value     float64   `json:"value"`
	Actual    *float64  `json:"actual"`
}

type ForecastResponse struct {
	RemainingCalls int        `json:"remainingCalls"`
	Forecasts      []Forecast `json:"forecasts"`
}
