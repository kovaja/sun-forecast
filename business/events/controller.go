package events

import (
	"database/sql"
	"fmt"
	"kovaja/sun-forecast/utils/logger"
	"time"
)

type AppEvent struct {
	Id        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

type EventController struct {
	repository EventRepository
}

func (ctl EventController) LogEvent(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	logger.Log("Event: %s", msg)

	err := ctl.repository.CreateEvent(msg)
	if err != nil {
		logger.LogError("Failed to write event", err)
	}
}

func (ctl EventController) ReadEvents() (*[]AppEvent, error) {
	return ctl.repository.ReadEvents()
}

func InitializeController(db *sql.DB) EventController {
	repository := EventRepository{
		db: db,
	}
	return EventController{
		repository: repository,
	}
}
