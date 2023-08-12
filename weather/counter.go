package weather

import (
	"kovaja/sun-forecast/logger"
	"kovaja/sun-forecast/utils"
	"time"
)

type RemainingCalls map[string]int

const MAX_CALLS = 1000
const CALLS_FILE = "./tmp/remaining-calls.json"

func CanCall() (bool, int) {
	remainingCalls, err := getRemainingCalls()

	if err != nil {
		logger.LogError("Failed to read remaining calls", err)
		return false, 0
	}

	newRemainingCalls := remainingCalls - 1
	canCall := newRemainingCalls > 0
	if canCall {
		err := setRemainigCalls(newRemainingCalls)

		if err != nil {
			logger.LogError("Failed to write remaining calls", err)
			return false, 0
		}
	}

	return canCall, newRemainingCalls
}

func getRemainingCalls() (int, error) {
	var data RemainingCalls
	err := utils.ReadJson(CALLS_FILE, &data)
	if err != nil {
		return 0, err
	}

	todayKey := getToday()
	remainingCalls, exists := data[todayKey]

	if !exists {
		addCurrentDay()
		logger.Log("Add new day for remaining calls")
		return MAX_CALLS, nil
	}

	return remainingCalls, nil
}

func addCurrentDay() error {
	return setRemainigCalls(MAX_CALLS)
}

func getToday() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02")
}

func setRemainigCalls(remainingCalls int) error {
	var data RemainingCalls
	err := utils.ReadJson(CALLS_FILE, &data)
	if err != nil {
		return err
	}

	todayKey := getToday()
	data[todayKey] = remainingCalls
	logger.Log("Setting remaining calls for %s as %d", todayKey, remainingCalls)

	return utils.WriteJson(CALLS_FILE, data)
}
