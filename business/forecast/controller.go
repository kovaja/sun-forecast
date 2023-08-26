package forecast

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"kovaja/sun-forecast/business/events"
	"kovaja/sun-forecast/utils"
	"kovaja/sun-forecast/utils/logger"
	"net/http"
	"time"
)

type ForecastController struct {
	repository *ForecastRepository
	eventCtl   *events.EventController
}

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

func (ctl ForecastController) readForecastsFromApi() (*ForecastResponse, error) {
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
			ctl.eventCtl.LogEvent("Failed to construct forecast %v. With error: %v", apiForecast, err)
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

func (ctl ForecastController) ConsumeForecasts() error {
	data, err := ctl.readForecastsFromApi()
	if err != nil {
		return err
	}

	added := 0
	updated := 0
	skipped := 0

	for _, forecast := range data.Forecasts {
		isExisting, id, value := ctl.repository.IsExistingForecastByPeriodEnd(forecast.PeriodEnd)

		if isExisting {
			if value != forecast.Value {
				err := ctl.repository.UpdateForcastValue(id, forecast.Value)
				if err != nil {
					return err
				}
				updated += 1
			} else {
				skipped += 1
			}
		} else {
			err := ctl.repository.createForcast(&forecast)
			if err != nil {
				return err
			}
			added += 1
		}
	}

	ctl.eventCtl.LogEvent("Updated forecasts, %d added, %d updated, %d skipped", added, updated, skipped)
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

func (ctl ForecastController) GetForecasts(fromStr string, toStr string) (*[]Forecast, error) {
	logger.Log("Read forecasts from DB: from %s to %s", fromStr, toStr)

	from, to, err := parseReadQuery(fromStr, toStr)
	if err != nil {
		return nil, utils.CustomError("Failed to parse query params", err)
	}

	return ctl.repository.ReadForecasts(from, to)
}

func (ctl ForecastController) UpdateForecasts(r *http.Request) ([]ForecastUpdate, error) {
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

	updates := ComputeUpdates(ctl.repository.GetExistingForecastByPeriodEnd, records)
	updated := 0

	for _, update := range updates {
		err := ctl.repository.UpdateForcastActual(&update)
		if err != nil {
			return nil, err
		}
		updated += 1
	}

	logger.Log("Updated forecasts with actual values, %d updated", updated)
	return updates, nil
}

func InitializeController(db *sql.DB, eventController events.EventController) ForecastController {
	repository := ForecastRepository{
		db: db,
	}

	return ForecastController{
		repository: &repository,
		eventCtl:   &eventController,
	}
}
