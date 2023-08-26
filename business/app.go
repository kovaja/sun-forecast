package app

import (
	"database/sql"
	"kovaja/sun-forecast/business/counter"
	"kovaja/sun-forecast/business/events"
	"kovaja/sun-forecast/business/forecast"
)

type AppControllers struct {
	EventCtl    events.EventController
	ForecastCtl forecast.ForecastController
}

func initializeEventController(db *sql.DB) events.EventController {
	repository := events.InitializeRepository(db)
	return events.InitializeController(repository)
}

func initializeForecastController(db *sql.DB, eventCtl events.EventController) forecast.ForecastController {
	repository := forecast.InitializeRepository(db)
	counterRepository := counter.InitializeRemainigCallRepository(db)
	counter := counter.InitializeCounter(counterRepository)

	return forecast.InitializeController(repository, counter, eventCtl)
}

func InitializeApp(db *sql.DB) AppControllers {
	eventCtl := initializeEventController(db)

	return AppControllers{
		EventCtl:    eventCtl,
		ForecastCtl: initializeForecastController(db, eventCtl),
	}
}
