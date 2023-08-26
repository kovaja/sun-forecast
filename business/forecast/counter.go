package forecast

import (
	"kovaja/sun-forecast/utils/logger"
	"time"
)

type RemainingCalls map[string]int

type Counter struct {
	repository RemainingCallRepository
}

func (c Counter) CanCall() (bool, int) {
	remainingCalls, err := c.getRemainingCalls()

	if err != nil {
		logger.LogError("Failed to read remaining calls", err)
		return false, 0
	}

	newRemainingCalls := remainingCalls - 1
	canCall := newRemainingCalls > 0
	if canCall {
		err := c.setRemainigCalls(newRemainingCalls)

		if err != nil {
			logger.LogError("Failed to write remaining calls", err)
			return false, 0
		}
	}

	return canCall, newRemainingCalls
}

func (c Counter) getRemainingCalls() (int, error) {
	todayKey := getToday()
	remaining, err := c.repository.ReadRemainingCalls(todayKey)

	logger.Log("Read remaining for %s from DB: %d", todayKey, remaining)
	return remaining, err
}

func (c Counter) setRemainigCalls(remainingCalls int) error {
	todayKey := getToday()
	return c.repository.UpdateRemainingCalls(todayKey, remainingCalls)
}

func getToday() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02")
}

func InitializeCounter(repo RemainingCallRepository) Counter {
	return Counter{
		repository: repo,
	}
}
