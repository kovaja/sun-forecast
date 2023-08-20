package events

import (
	"fmt"
	"kovaja/sun-forecast/db"
	"kovaja/sun-forecast/logger"
	"kovaja/sun-forecast/utils"
	"time"
)

type AppEvent struct {
	Id        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

func LogEvent(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	logger.Log("Event: %s", msg)

	db := db.GetDb()
	query := "INSERT INTO events (message) VALUES ($1)"
	_, err := db.Exec(query, msg)

	if err != nil {
		logger.LogError("Failed to write event", err)
	}
}

func ReadEvents() (*[]AppEvent, error) {
	db := db.GetDb()
	query := "SELECT id, timestamp, message FROM events ORDER BY timestamp DESC;"

	rows, err := db.Query(query)
	if err != nil {
		return nil, utils.CustomError("Failed to read events", err)
	}
	defer rows.Close()

	var events []AppEvent
	for rows.Next() {
		var event AppEvent
		err := rows.Scan(&event.Id, &event.Timestamp, &event.Message)
		if err != nil {
			return nil, utils.CustomError("Failed to read single event", err)
		}

		events = append(events, event)
	}

	return &events, nil
}
