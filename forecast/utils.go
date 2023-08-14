package forecast

import (
	"kovaja/sun-forecast/logger"
	"strconv"
	"time"
)

func getPeriodEnd(record HaHistoryRecord) time.Time {
	inputTime := record.LastChanged
	fullHourDiff := 60 - inputTime.Minute()

	// logger.Log("Process input time %v, fullhourdiff %d", inputTime, fullHourDiff)
	var updatedTime time.Time
	if fullHourDiff < 30 {
		// there is less than 30 minutes till full hour add the rest to the time
		// logger.Log("Adding time till full hour %d]",fullHourDiff)
		updatedTime = inputTime.Add(time.Minute * time.Duration(fullHourDiff))
	} else {
		// there is more than 30 minutes till full hour, add only the rest till half hour
		// 43 - 30 = 13 | 13.17 => 13.30
		// logger.Log("Adding time till half hour %d]",fullHourDiff-30)
		updatedTime = inputTime.Add(time.Minute * time.Duration(fullHourDiff-30))
	}

	result := time.Date(
		updatedTime.Year(),
		updatedTime.Month(),
		updatedTime.Day(),
		updatedTime.Hour(),
		updatedTime.Minute(),
		0, // seconds and nanoseconds don't matter
		0,
		updatedTime.Location(),
	)

	// logger.Log("Period end for %v computed as %v", inputTime, result)

	return result
}

func ComputeUpdates(records []HaHistoryRecord) []ForecastUpdate {
	/**
		input record:
		{
			"state": "505",
			"last_changed": "2023-08-14T16:16:29.758217+00:00"
	  }
		need to recognize period end - that is next half an hour
		need to convert state string to float64
	*/

	var updates []ForecastUpdate
	var currentPeriodEnd *time.Time = nil
	periodSum := 0.0
	periodLen := 0

	for _, record := range records {
		periodEnd := getPeriodEnd(record)
		actual, err := strconv.ParseFloat(record.State, 64)

		if err != nil {
			logger.LogError("Failed to convert record state", err)
		} else {
			if currentPeriodEnd == nil || currentPeriodEnd.UnixMicro() != periodEnd.UnixMicro() {
				if currentPeriodEnd != nil {
					// logger.Log("Adding update for %v, len: %d, sum %.2f, avg: %.2f", currentPeriodEnd, periodLen, periodSum, periodSum/float64(periodLen))
					updates = append(updates, ForecastUpdate{PeriodEnd: *currentPeriodEnd, Actual: periodSum / float64(periodLen)})
				}
				currentPeriodEnd = &periodEnd
				periodSum = 0.0
				periodLen = 0
			}

			periodSum += actual
			periodLen += 1
		}
	}

	// add the last one period
	// logger.Log("Adding update for %v, len: %d, sum %.2f, avg: %.2f", currentPeriodEnd, periodLen, periodSum, periodSum/float64(periodLen))
	updates = append(updates, ForecastUpdate{PeriodEnd: *currentPeriodEnd, Actual: periodSum / float64(periodLen)})

	logger.Log("Computed updated: %v", updates)

	return updates
}
