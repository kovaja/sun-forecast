package forecast

import (
	"database/sql"
	"fmt"
	"kovaja/sun-forecast/utils"
	"kovaja/sun-forecast/utils/logger"
	"time"
)

type ForecastRepository struct {
	db *sql.DB
}

func (repository ForecastRepository) GetExistingForecastByPeriodEnd(timestamp time.Time) *Forecast {
	query := "SELECT * FROM forecasts WHERE period_end = $1"

	var forecast Forecast
	err := repository.db.QueryRow(query, timestamp).Scan(&forecast.Id, &forecast.PeriodEnd, &forecast.Value, &forecast.Actual, &forecast.ActualCount, &forecast.LastActualAt)
	if err == sql.ErrNoRows {
		logger.LogError("Did not find existing actuals", err)
		return nil
	} else if err != nil {
		logger.LogError("Failed to load existing forecast", err)
		return nil
	}

	return &forecast
}

func (repository ForecastRepository) IsExistingForecastByPeriodEnd(timestamp time.Time) (bool, int, float64) {
	forecast := repository.GetExistingForecastByPeriodEnd(timestamp)

	if forecast == nil {
		return false, -1, -1
	}

	return true, forecast.Id, forecast.Value
}

func (repository ForecastRepository) UpdateForcastActual(update *ForecastUpdate) error {
	query := "UPDATE forecasts SET actual = $1, actual_count = $2, last_actual_at = $3 WHERE period_end = $4"
	result, err := repository.db.Exec(query, update.Actual, update.ActualCount, update.LastActualAt, update.PeriodEnd)

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

func (repository ForecastRepository) ReadForecasts(from *time.Time, to *time.Time) (*[]Forecast, error) {
	query := "SELECT id, period_end, value, actual, actual_count, last_actual_at FROM forecasts WHERE period_end >= $1 AND period_end <= $2 ORDER BY period_end ASC"

	rows, err := repository.db.Query(query, from, to)
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

func (repository ForecastRepository) UpdateForcastValue(id int, value float64) error {
	query := "UPDATE forecasts SET value = $1 WHERE id = $2"
	result, err := repository.db.Exec(query, value, id)

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

func (repository ForecastRepository) createForcast(forecast *Forecast) error {
	query := "INSERT INTO forecasts (period_end, value) VALUES ($1, $2)"

	_, err := repository.db.Exec(query, forecast.PeriodEnd, forecast.Value)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to insert forecast %v. Error: %v", forecast, err)
		return utils.CustomError(errorMsg, err)
	}

	return nil
}

func InitializeRepository(db *sql.DB) ForecastRepository {
	return ForecastRepository{
		db: db,
	}
}
