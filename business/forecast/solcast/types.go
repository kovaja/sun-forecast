package solcast

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
