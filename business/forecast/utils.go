package forecast

import (
	"kovaja/sun-forecast/business/forecast/solcast"
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
		// logger.Log("Adding time till full hour %d",fullHourDiff)
		updatedTime = inputTime.Add(time.Minute * time.Duration(minuteDiff))
	} else {
		// there is more than 30 minutes till full hour, add only the rest till half hour
		// 43 - 30 = 13 | 13.17 => 13.30
		// logger.Log("Adding time till half hour %d",fullHourDiff-30)
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

func appendUpdate(updates []ForecastUpdate, updated bool, newUpdate ForecastUpdate) []ForecastUpdate {
	// add the last one period
	if !updated {
		logger.Log("Will not add %v between updates as no change was made", newUpdate.PeriodEnd)
		return updates
	}

	logger.Log("Adding new update %v", newUpdate)

	return append(updates, newUpdate)
}

/*
	  input record:
		{
			"state": "505",
			"last_changed": "2023-08-14T16:16:29.758217+00:00"
		}
		need to recognize period end - that is next half an hour
		need to convert state string to float64
*/
func ComputeUpdates(loadExistingForecast func(time.Time) *Forecast, records []HaHistoryRecord) []ForecastUpdate {
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
			skippedRecords += 1
		} else {
			if currentPeriodEnd == nil || currentPeriodEnd.UnixMicro() != periodEnd.UnixMicro() {
				if currentPeriodEnd != nil {
					updates = appendUpdate(
						updates,
						updated,
						ForecastUpdate{
							PeriodEnd:    *currentPeriodEnd,
							Actual:       average,
							ActualCount:  count,
							LastActualAt: *lastChanged,
						},
					)
				}

				currentPeriodEnd = &periodEnd
				existingForecast := loadExistingForecast(*currentPeriodEnd)
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
				// this should rather be weighted average reflecting the period of time this state lasted
				// but so far it seems the difference is not that dramatic
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

	if updated {
		updates = appendUpdate(
			updates,
			updated,
			ForecastUpdate{
				PeriodEnd:    *currentPeriodEnd,
				Actual:       average,
				ActualCount:  count,
				LastActualAt: *lastChanged,
			},
		)
	}

	logger.Log("Computed updates for %d periods, %d records skipped", len(updates), skippedRecords)
	return updates
}

/*
	  For given collection of forecasts it returns time from and to.
		From is start of period of first forecast
		To is end of period of last forecast
*/
func GetForecastsRange(forecasts []Forecast) (*time.Time, *time.Time) {
	if len(forecasts) == 0 {
		// we have no forecasts return nil
		return nil, nil
	}

	firstForecast := forecasts[0]
	// subtract 30 minutes to get us to period start
	periodStart := firstForecast.PeriodEnd.Add(time.Minute * -time.Duration(30))

	if len(forecasts) == 1 {
		return &periodStart, &firstForecast.PeriodEnd
	}

	lastForecast := forecasts[len(forecasts)-1]

	return &periodStart, &lastForecast.PeriodEnd
}

func ParseReadQuery(fromStr string, toStr string) (*time.Time, *time.Time, error) {
	fromQueryTime, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return nil, nil, err
	}

	toQueryTime, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return nil, nil, err
	}

	// from will stay the same, period_end >= from will cover that time
	// it will be the first half an hour after requested time
	from := fromQueryTime
	// for 16:15 we need add 30 minutes to get 16:45, that would return 16:30 period_end
	// then we can query period_end <= to
	to := toQueryTime.Add(time.Minute * time.Duration(30))

	return &from, &to, nil
}

func ConstructForcast(apiForecast solcast.SolcastApiForecast) (*Forecast, error) {
	parsedTime, err := time.Parse(TS_LAYOUT, apiForecast.PeriodEnd)
	if err != nil {
		return nil, err
	}

	forecast := Forecast{
		Id:        -1,
		PeriodEnd: parsedTime,
		Value:     apiForecast.PvEstimate * 1000, // convert 1.56kw na 1560w
		Actual:    nil,
	}

	return &forecast, nil
}
