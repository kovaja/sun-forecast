package events

import (
	"fmt"
	"kovaja/sun-forecast/utils/logger"
)

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

func (ctl EventController) ReadEvents(eventType string) (*[]AppEvent, error) {
	if eventType == "" {
		return ctl.repository.ReadEvents()
	}

	requestedEventType := getValidEventTypeQueryParam(eventType)
	return ctl.repository.ReadEventsByType(int(requestedEventType))
}

func InitializeController(repository EventRepository) EventController {
	return EventController{
		repository: repository,
	}
}
