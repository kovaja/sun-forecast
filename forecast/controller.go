package forecast

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"kovaja/sun-forecast/db"
	"kovaja/sun-forecast/events"
	"kovaja/sun-forecast/logger"
	"kovaja/sun-forecast/utils"
	"net/http"
	"time"
)

const tsLayout = "2006-01-02T15:04:05.0000000Z"

func constructForcast(apiForecast SolcastApiForecast) (*Forecast, error) {
	parsedTime, err := time.Parse(tsLayout, apiForecast.PeriodEnd)
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
		return nil, errors.New("Cannot call solcast api (Too many calls today)")
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

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to check affected rows/ for forecast %d", id)
		return utils.CustomError(errorMsg, err)
	}

	logger.Log("Updated forecast value %d, rowsAffected: %d", id, rowsAffected)
	return nil
}

func updateForcastActual(db *sql.DB, periodEnd time.Time, actual float64) error {
	query := "UPDATE forecasts SET actual = $1 WHERE period_end = $2"
	result, err := db.Exec(query, actual, periodEnd)

	if err != nil {
		errorMsg := fmt.Sprintf("Failed to update forecast %v.", periodEnd)
		return utils.CustomError(errorMsg, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to check affected rows/ for forecast %v", periodEnd)
		return utils.CustomError(errorMsg, err)
	}

	logger.Log("Updated forecast actual %v, rowsAffected: %d", periodEnd, rowsAffected)
	return nil
}

func createForcast(db *sql.DB, forecast *Forecast) error {
	query := "INSERT INTO forecasts (period_end, value) VALUES ($1, $2)"

	_, err := db.Exec(query, forecast.PeriodEnd, forecast.Value)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to insert forecast %v. Error: %v", forecast, err)
		return utils.CustomError(errorMsg, err)
	}

	logger.Log("Created forecast for %s", forecast.PeriodEnd)
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
				logger.Log("Skipping forecast updated for %d as value is the same", id)
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

func ReadForecastsFromDb() (*[]Forecast, error) {
	db := db.GetDb()
	query := "SELECT id, period_end, value, actual FROM forecasts"

	rows, err := db.Query(query)
	if err != nil {
		return nil, utils.CustomError("Failed to read singe forecast", err)
	}
	defer rows.Close()

	var forecasts []Forecast
	for rows.Next() {
		var forecast Forecast
		err := rows.Scan(&forecast.Id, &forecast.PeriodEnd, &forecast.Value, &forecast.Actual)
		if err != nil {
			return nil, utils.CustomError("Failed to read singe forecast", err)
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

	updates := ComputeUpdates(records)
	db := db.GetDb()
	updated := 0

	for _, update := range updates {
		err := updateForcastActual(db, update.PeriodEnd, update.Actual)
		if err != nil {
			return nil, err
		}
		updated += 1
	}

	events.LogEvent("Updated forecasts with actual values, %d updated", updated)
	return updates, nil
}
