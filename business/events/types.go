package events

import "time"

type EventType int

const (
	ForecastConsumed EventType = iota
	ForecastUpdated
	AppError
	// new event goes here
	NoEvents
)

type AppEvent struct {
	Id        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Type      EventType `json:"type"`
}
