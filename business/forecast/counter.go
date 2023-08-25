package forecast

import (
	"database/sql"
	"kovaja/sun-forecast/core/db"
	"kovaja/sun-forecast/utils/logger"
	"time"
)

type RemainingCalls map[string]int

const MAX_CALLS = 8

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
	todayKey := getToday()
	db := db.GetDb()

	var remaining int
	err := db.QueryRow("SELECT remaining FROM remaining_calls WHERE date = $1", todayKey).Scan(&remaining)

	if err != nil {
		if err == sql.ErrNoRows {
			// this means that we don't have that day yet, so return MAX_CALLS
			// we will write it later on
			logger.Log("No record found for %s, returning MAX_CALLS", todayKey)
			remaining = MAX_CALLS
		} else {
			return 0, err
		}
	}

	logger.Log("Read remaining for %s from DB: %d", todayKey, remaining)
	return remaining, nil
}

func getToday() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02")
}

func setRemainigCalls(remainingCalls int) error {
	todayKey := getToday()
	db := db.GetDb()
	query := `
    INSERT INTO remaining_calls (date, remaining)
    VALUES ($1, $2)
    ON CONFLICT (date) DO UPDATE SET remaining = EXCLUDED.remaining
  `
	_, err := db.Exec(query, todayKey, remainingCalls)

	logger.Log("Setting remaining calls for %s as %d", todayKey, remainingCalls)
	return err
}
