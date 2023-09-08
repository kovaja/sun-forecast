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

func getValidLimit(limitParam string) int {
	num, err := strconv.Atoi(limitParam)
	if err != nil {
		return 1000000000
	}
	return num
}
