package forecast

import (
	"encoding/json"
	"errors"
	"fmt"
	"kovaja/sun-forecast/business/counter"
	"kovaja/sun-forecast/business/events"
	"kovaja/sun-forecast/business/forecast/solcast"
	"kovaja/sun-forecast/utils"
	"kovaja/sun-forecast/utils/logger"
	"net/http"
)

type ForecastController struct {
	repository *ForecastRepository
	counter    *counter.Counter
	eventCtl   *events.EventController
}

const TS_LAYOUT = "2006-01-02T15:04:05.0000000Z"

var (
	ErrSolcastTooManyCalls = errors.New("Cannot call solcast api (Too many calls today)")
)

func (ctl ForecastController) readForecastsFromApi() ([]Forecast, error) {
	canCall, remainingCalls := ctl.counter.CanCall()

	if !canCall {
		return nil, ErrSolcastTooManyCalls
	}

	body, err := solcast.FetchForecasts()
	if err != nil {
		return nil, err
	}

	var forcasts []Forecast
	for _, apiForecast := range body.Forecasts {
		forecast, err := ConstructForcast(apiForecast)

		if err != nil {
			ctl.eventCtl.LogAppError("Failed to construct forecast %v. With error: %v", apiForecast, err)
		} else {
			forcasts = append(forcasts, *forecast)
		}
	}

	logger.Log("Read %d forecasts from api. Calls remaining: %d", len(forcasts), remainingCalls)
	return forcasts, nil
}

func (ctl ForecastController) ConsumeForecasts() error {
	data, err := ctl.readForecastsFromApi()
	if err != nil {
		return err
	}

	added := 0
	updated := 0
	skipped := 0

	for _, forecast := range data {
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

	ctl.eventCtl.LogEvent(events.ForecastConsumed, "Updated forecasts, %d added, %d updated, %d skipped", added, updated, skipped)
	return nil
}

func (ctl ForecastController) GetForecasts(fromStr string, toStr string) (*ForecastResponse, error) {
	logger.Log("Read forecasts from DB: from %s to %s", fromStr, toStr)

	from, to, err := ParseReadQuery(fromStr, toStr)
	if err != nil {
		return nil, utils.CustomError("Failed to parse query params", err)
	}

	forecasts, err := ctl.repository.ReadForecasts(from, to)
	if err != nil {
		return nil, err
	}

	fromStart, toEnd := GetForecastsRange(*forecasts)

	return &ForecastResponse{
		Forecasts: *forecasts,
		From:      *fromStart,
		To:        *toEnd,
	}, nil
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

	ctl.eventCtl.LogEvent(
		events.ForecastUpdated,
		"Updated forecasts with actual values, records %d, %d forecasts updated",
		len(records),
		updated,
	)
	return updates, nil
}

func InitializeController(
	repository ForecastRepository,
	counter counter.Counter,
	eventController events.EventController,
) ForecastController {

	return ForecastController{
		repository: &repository,
		eventCtl:   &eventController,
		counter:    &counter,
	}
}
