package events

import (
	"fmt"
	"kovaja/sun-forecast/db"
	"kovaja/sun-forecast/logger"
)

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
