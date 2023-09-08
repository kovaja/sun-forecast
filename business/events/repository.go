package events

import (
	"database/sql"
	"kovaja/sun-forecast/utils"
)

type EventRepository struct {
	db *sql.DB
}

func (repository EventRepository) CreateEvent(eventType int, message string) error {
	query := "INSERT INTO events (message, type) VALUES ($1, $2)"
	_, err := repository.db.Exec(query, message, eventType)
	return err
}

func processReadEvents(rows *sql.Rows, err error) (*[]AppEvent, error) {
	if err != nil {
		return nil, utils.CustomError("Failed to read events", err)
	}
	defer rows.Close()

	var events []AppEvent
	for rows.Next() {
		var event AppEvent
		err := rows.Scan(&event.Id, &event.Timestamp, &event.Message, &event.Type)
		if err != nil {
			return nil, utils.CustomError("Failed to read single event", err)
		}

		events = append(events, event)
	}

	return &events, nil
}

func (repository EventRepository) ReadEvents(limit int) (*[]AppEvent, error) {
	query := "SELECT id, timestamp, message, type FROM events ORDER BY timestamp DESC LIMIT $1;"
	rows, err := repository.db.Query(query, limit)
	return processReadEvents(rows, err)
}

func (repository EventRepository) ReadEventsByType(eventType int, limit int) (*[]AppEvent, error) {
	query := "SELECT id, timestamp, message, type FROM events WHERE type = $1 ORDER BY timestamp DESC LIMIT $2;"
	rows, err := repository.db.Query(query, eventType, limit)
	return processReadEvents(rows, err)
}

func InitializeRepository(db *sql.DB) EventRepository {
	return EventRepository{
		db: db,
	}
}
