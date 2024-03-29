package forecast

import "time"

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

type ForecastDiff struct {
	Date time.Time `json:"date"`
	// the length of this array is always the same
	// we have 30minutes invervals, therefore there will be 48 diffs
	Diffs []float64 `json:"diffs"`
}
