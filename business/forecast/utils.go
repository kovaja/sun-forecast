package forecast

import (
	"database/sql"
	"kovaja/sun-forecast/utils/logger"
	"strconv"
	"time"
)

func getPeriodEnd(record HaHistoryRecord) time.Time {
	inputTime := record.LastChanged
	minuteDiff := 60 - inputTime.Minute()

	// logger.Log("Process input time %v, fullhourdiff %d", inputTime, fullHourDiff)
	var updatedTime time.Time
	if minuteDiff <= 30 {
		// there is less than 30 minutes till full hour add the rest to the time
		// logger.Log("Adding time till full hour %d]",fullHourDiff)
		updatedTime = inputTime.Add(time.Minute * time.Duration(minuteDiff))
	} else {
		// there is more than 30 minutes till full hour, add only the rest till half hour
		// 43 - 30 = 13 | 13.17 => 13.30
		// logger.Log("Adding time till half hour %d]",fullHourDiff-30)
		updatedTime = inputTime.Add(time.Minute * time.Duration(minuteDiff-30))
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

func appendUpdate(
	updates []ForecastUpdate,
	updated bool,
	currentPeriodEnd *time.Time,
	average float64,
	count int,
	lastChanged *time.Time,
) []ForecastUpdate {
	// add the last one period
	if !updated {
		logger.Log("Will not add %v between updates as no change was made", currentPeriodEnd)
		return updates
	}

	newUpdate := ForecastUpdate{
		PeriodEnd:    *currentPeriodEnd,
		Actual:       average,
		ActualCount:  count,
		LastActualAt: *lastChanged,
	}

	logger.Log("Adding new update %v", newUpdate)

	return append(updates, newUpdate)
}

func ComputeUpdates(db *sql.DB, records []HaHistoryRecord) []ForecastUpdate {
	/**
		input record:
		{
			"state": "505",
			"last_changed": "2023-08-14T16:16:29.758217+00:00"
	  }
		need to recognize period end - that is next half an hour
		need to convert state string to float64
	*/

	skippedRecords := 0
	var updates []ForecastUpdate
	var currentPeriodEnd *time.Time = nil
	average := 0.0
	count := 0
	var lastChanged *time.Time = nil
	updated := false

	for _, record := range records {
		periodEnd := getPeriodEnd(record)
		actual, err := strconv.ParseFloat(record.State, 64)
		recordLastChanged := record.LastChanged

		if err != nil {
			logger.LogError("Failed to convert actual", err)
		} else {
			if currentPeriodEnd == nil || currentPeriodEnd.UnixMicro() != periodEnd.UnixMicro() {
				if currentPeriodEnd != nil {
					updates = appendUpdate(updates, updated, currentPeriodEnd, average, count, lastChanged)
				}

				currentPeriodEnd = &periodEnd
				existingForecast := getExistingForecast(db, *currentPeriodEnd)
				updated = false
				if existingForecast != nil {
					lastChanged = existingForecast.LastActualAt
					if lastChanged != nil {
						average = *existingForecast.Actual
						count = *&existingForecast.ActualCount
						logger.Log("lastchange for %v not null thus setup %.2f, %d, %v", currentPeriodEnd, average, count, lastChanged)
					} else {
						// this is foreast that has never been touched
						logger.Log("lastchange for %v is null therfore setting things to zero", currentPeriodEnd)
						average = 0
						count = 0
					}
				} else {
					// this shoulnd't really happend as we should always have forecast to update
					logger.Log("No existing forecast for %v", periodEnd)
					lastChanged = nil
					average = 0
					count = 0
				}
			}

			if lastChanged == nil || record.LastChanged.After(*lastChanged) {
				// cumulative average for every record value
				average = (average*float64(count) + actual) / (float64(count + 1))
				count += 1
				lastChanged = &recordLastChanged
				updated = true
			} else {
				logger.Log("Skipping record as average already contains this one %v, record last change %v", lastChanged, record.LastChanged)
				skippedRecords += 1
			}
		}
	}

	updates = appendUpdate(updates, updated, currentPeriodEnd, average, count, lastChanged)
	logger.Log("Computed updates for %d periods, %d records skipped", len(updates), skippedRecords)
	return updates
}
