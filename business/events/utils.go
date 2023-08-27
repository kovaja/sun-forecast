package events

import (
	"strconv"
)

func getValidEventTypeQueryParam(eventTypeParam string) EventType {
	num, err := strconv.Atoi(eventTypeParam)
	if err != nil {
		return NoEvents
	}

	switch num {
	case 0:
		return ForecastConsumed
	case 1:
		return ForecastUpdated
	case 2:
		return AppError
	}

	return NoEvents
}
