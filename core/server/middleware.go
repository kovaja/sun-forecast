package server

import (
	"kovaja/sun-forecast/utils/logger"
	"net/http"
	"time"
)

func logRequest(handler ApiHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		logger.Log("%s %s", r.Method, r.URL)
		data, err := handler(r)
		SendResponse(w, data, err)

		duration := time.Since(startTime)
		logger.Log("%s %s took %dms", r.Method, r.URL, duration.Milliseconds())
	}
}
