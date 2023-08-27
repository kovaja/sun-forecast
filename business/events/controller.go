package events

import (
	"fmt"
	"kovaja/sun-forecast/utils/logger"
	"time"
)

type EventType int

const (
	ForecastConsumed EventType = iota
	ForecastUpdated
	AppError
)

type AppEvent struct {
	Id        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Type      EventType `json:"type"`
}

type EventController struct {
	repository EventRepository
}

func (ctl EventController) LogEvent(eventType EventType, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	logger.Log("Event: %s", msg)

	err := ctl.repository.CreateEvent(int(eventType), msg)
	if err != nil {
		logger.LogError("Failed to write event", err)
	}
}

func (ctl EventController) LogAppError(format string, a ...any) {
	ctl.LogEvent(AppError, format, a...)
}

func (ctl EventController) ReadEvents() (*[]AppEvent, error) {
	return ctl.repository.ReadEvents()
}

func InitializeController(repository EventRepository) EventController {
	return EventController{
		repository: repository,
	}
}
