package forecast

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"kovaja/sun-forecast/business/events"
	"kovaja/sun-forecast/core/db"
	"kovaja/sun-forecast/utils"
	"kovaja/sun-forecast/utils/logger"
	"net/http"
	"time"
)

const TS_LAYOUT = "2006-01-02T15:04:05.0000000Z"

var (
	ErrSolcastTooManyCalls = errors.New("Cannot call solcast api (Too many calls today)")
)

func constructForcast(apiForecast SolcastApiForecast) (*Forecast, error) {
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

func readForecastsFromApi() (*ForecastResponse, error) {
	canCall, remainingCalls := CanCall()

	if !canCall {
		return nil, ErrSolcastTooManyCalls
	}

	body, err := fetchForecasts()
	if err != nil {
		return nil, err
	}

	var forcasts []Forecast
	for _, apiForecast := range body.Forecasts {
		forecast, err := constructForcast(apiForecast)

		if err != nil {
			events.LogEvent("Failed to construct forecast %v. With error: %v", apiForecast, err)
		} else {
			forcasts = append(forcasts, *forecast)
		}
	}

	response := ForecastResponse{
		RemainingCalls: remainingCalls,
		Forecasts:      forcasts,
	}

	return &response, nil
}

func isExistingForecast(db *sql.DB, timestamp time.Time) (bool, int, float64) {
	query := "SELECT id, value FROM forecasts WHERE period_end = $1"
	var id int
	var value float64
	err := db.QueryRow(query, timestamp).Scan(&id, &value)
	if err == sql.ErrNoRows {
		return false, -1, -1
	} else if err != nil {
		logger.LogError("Failed to check forcast record", err)

		// returning false to prevent duplicate records in the DB
		// rather not store anything that have multiple records for same timestamp
		return false, -1, -1
	}

	return true, id, value
}

func updateForcastValue(db *sql.DB, id int, value float64) error {
	query := "UPDATE forecasts SET value = $1 WHERE id = $2"
	result, err := db.Exec(query, value, id)

	if err != nil {
		errorMsg := fmt.Sprintf("Failed to update forecast %d.", id)
		return utils.CustomError(errorMsg, err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to check affected rows/ for forecast %d", id)
		return utils.CustomError(errorMsg, err)
	}

	return nil
}

func updateForcastActual(db *sql.DB, update *ForecastUpdate) error {
	query := "UPDATE forecasts SET actual = $1, actual_count = $2, last_actual_at = $3 WHERE period_end = $4"
	result, err := db.Exec(query, update.Actual, update.ActualCount, update.LastActualAt, update.PeriodEnd)

	if err != nil {
		errorMsg := fmt.Sprintf("Failed to update forecast %v.", update.PeriodEnd)
		return utils.CustomError(errorMsg, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to check affected rows/ for forecast %v", update.PeriodEnd)
		return utils.CustomError(errorMsg, err)
	}

	logger.Log("Updated forecast actual %v, rowsAffected: %d", update.PeriodEnd, rowsAffected)
	return nil
}

func createForcast(db *sql.DB, forecast *Forecast) error {
	query := "INSERT INTO forecasts (period_end, value) VALUES ($1, $2)"

	_, err := db.Exec(query, forecast.PeriodEnd, forecast.Value)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to insert forecast %v. Error: %v", forecast, err)
		return utils.CustomError(errorMsg, err)
	}

	return nil
}

func ConsumeForecasts() error {
	data, err := readForecastsFromApi()
	if err != nil {
		return err
	}

	db := db.GetDb()
	added := 0
	updated := 0
	skipped := 0

	for _, forecast := range data.Forecasts {
		isExisting, id, value := isExistingForecast(db, forecast.PeriodEnd)

		if isExisting {
			if value != forecast.Value {
				err := updateForcastValue(db, id, forecast.Value)
				if err != nil {
					return err
				}
				updated += 1
			} else {
				skipped += 1
			}
		} else {
			err := createForcast(db, &forecast)
			if err != nil {
				return err
			}
			added += 1
		}
	}

	events.LogEvent("Updated forecasts, %d added, %d updated, %d skipped", added, updated, skipped)
	return nil
}

func parseReadQuery(fromStr string, toStr string) (*time.Time, *time.Time, error) {
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

func ReadForecastsFromDb(fromStr string, toStr string) (*[]Forecast, error) {
	logger.Log("Read forecasts from DB: from %s to %s", fromStr, toStr)

	from, to, err := parseReadQuery(fromStr, toStr)
	if err != nil {
		return nil, utils.CustomError("Failed to parse query params", err)
	}

	db := db.GetDb()
	query := "SELECT id, period_end, value, actual, actual_count, last_actual_at FROM forecasts WHERE period_end >= $1 AND period_end <= $2 ORDER BY period_end ASC"

	rows, err := db.Query(query, from, to)
	if err != nil {
		return nil, utils.CustomError("Failed to read singe forecast", err)
	}
	defer rows.Close()

	var forecasts []Forecast
	for rows.Next() {
		var forecast Forecast
		err := rows.Scan(&forecast.Id, &forecast.PeriodEnd, &forecast.Value, &forecast.Actual, &forecast.ActualCount, &forecast.LastActualAt)
		if err != nil {
			return nil, utils.CustomError("Failed to read single forecast", err)
		}

		forecasts = append(forecasts, forecast)
	}

	return &forecasts, nil
}

func UpdateForecasts(r *http.Request) ([]ForecastUpdate, error) {
	var data [][]HaHistoryRecord
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, utils.CustomError("Failed to decode data for forecast update", err)
	}

	if len(data) == 0 || len(data) > 1 {
		return nil, errors.New(fmt.Sprintf("Unexpected length of first data array: %d", len(data)))
	}

	records := data[0]
	logger.Log("Received update data %d", len(records))

	db := db.GetDb()
	updates := ComputeUpdates(getExistingForecastLoader(db), records)
	updated := 0

	for _, update := range updates {
		err := updateForcastActual(db, &update)
		if err != nil {
			return nil, err
		}
		updated += 1
	}

	logger.Log("Updated forecasts with actual values, %d updated", updated)
	return updates, nil
}
