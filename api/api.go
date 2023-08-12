package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type Response map[string]any
type ErrorResponse struct {
	Error string `json:"error"`
}

func sendJson(w http.ResponseWriter, responseData Response, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

func getResponse(data any, err error) (Response, error) {
	response := make(Response)

	response["date"] = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	response["data"] = data

	return response, err
}

func SendResponse(w http.ResponseWriter, data any, err error) {
	response, err := getResponse(data, err)
	sendJson(w, response, err)
}
