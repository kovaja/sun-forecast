package events

import (
	"database/sql"
	"kovaja/sun-forecast/utils"
)

type EventRepository struct {
	db *sql.DB
}

func (repository EventRepository) CreateEvent(message string) error {
	query := "INSERT INTO events (message) VALUES ($1)"
	_, err := repository.db.Exec(query, message)
	return err
}

func (repository EventRepository) ReadEvents() (*[]AppEvent, error) {
	query := "SELECT id, timestamp, message FROM events ORDER BY timestamp DESC;"

	rows, err := repository.db.Query(query)
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

func InitializeRepository(db *sql.DB) EventRepository {
	return EventRepository{
		db: db,
	}
}
