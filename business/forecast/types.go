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
	Id           int        `json:"id"`
	PeriodEnd    time.Time  `json:"periodEnd"`
	Value        float64    `json:"value"`
	Actual       *float64   `json:"actual"`
	ActualCount  int        `json:"actualCount"`
	LastActualAt *time.Time `json:"lastActualAt"`
}

type ForecastResponse struct {
	Forecasts []Forecast `json:"forecasts"`
	From      time.Time  `json:"from"`
	To        time.Time  `json:"to"`
}

type HaHistoryRecord struct {
	LastChanged time.Time `json:"last_changed"`
	State       string    `json:"state"`
}

type ForecastUpdate struct {
	PeriodEnd    time.Time
	Actual       float64
	ActualCount  int
	LastActualAt time.Time
}
