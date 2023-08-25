package forecast

import (
	"database/sql"
	"kovaja/sun-forecast/utils/logger"
	"time"
)

func getExistingForecast(db *sql.DB, timestamp time.Time) *Forecast {
	query := "SELECT * FROM forecasts WHERE period_end = $1"

	var forecast Forecast
	err := db.QueryRow(query, timestamp).Scan(&forecast.Id, &forecast.PeriodEnd, &forecast.Value, &forecast.Actual, &forecast.ActualCount, &forecast.LastActualAt)
	if err == sql.ErrNoRows {
		logger.LogError("Did not find existing actuals", err)
		return nil
	} else if err != nil {
		logger.LogError("Failed to load existing forecast", err)
		return nil
	}

	return &forecast
}
