package events

import "database/sql"

type EventRepository struct {
	db *sql.DB
}
