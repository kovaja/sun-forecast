package events

import (
	"fmt"
	"kovaja/sun-forecast/utils/logger"
)

type EventController struct {
	repository EventRepository
}

func getMessage(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func logEventByLogger(msg string) {
	logger.Log("Event: %s", msg)
}

func (ctl EventController) LogEvent(eventType EventType, format string, a ...any) {
	msg := getMessage(format, a...)
	logEventByLogger(msg)

	err := ctl.repository.CreateEvent(int(eventType), msg)
	if err != nil {
		logger.LogError("Failed to write event", err)
	}
}

func (ctl EventController) LogEventIf(prereq bool, eventType EventType, format string, a ...any) {
	if prereq {
		ctl.LogEvent(eventType, format, a...)
	} else {
		// just log it
		logEventByLogger(getMessage(format, a...))
	}
}

func (ctl EventController) LogAppError(format string, a ...any) {
	ctl.LogEvent(AppError, format, a...)
}

func (ctl EventController) ReadEvents(eventType string, limitParam string) (*[]AppEvent, error) {
	limit := getValidLimit(limitParam)
	if eventType == "" {
		return ctl.repository.ReadEvents(limit)
	}

	requestedEventType := getValidEventTypeQueryParam(eventType)
	return ctl.repository.ReadEventsByType(int(requestedEventType), limit)
}

func InitializeController(repository EventRepository) EventController {
	return EventController{
		repository: repository,
	}
}
