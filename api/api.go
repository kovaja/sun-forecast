package api

import (
	"encoding/json"
	"net/http"
	"reflect"
	"time"
)

type Response map[string]any
type ErrorResponse struct {
	Error string `json:"error"`
}

func sendJson(w http.ResponseWriter, responseData Response, err error, status *int) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		finalStatus := http.StatusBadRequest
		if status != nil {
			finalStatus = *status
		}
		w.WriteHeader(finalStatus)

		errorRespone := ErrorResponse{Error: err.Error()}
		errorResponseJson, err := json.Marshal(errorRespone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(errorResponseJson)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}

func maybeSliceLen(data interface{}) int {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() == reflect.Slice {
		return value.Len()
	}

	return -1
}

func getResponse(data any) Response {
	response := make(Response)

	response["date"] = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	response["data"] = data
	response["num"] = maybeSliceLen(data)

	return response
}

func SendResponse(w http.ResponseWriter, data any, err error) {
	response := getResponse(data)
	sendJson(w, response, err, nil)
}

func SendError(w http.ResponseWriter, err error, status int) {
	sendJson(w, nil, err, &status)
}
